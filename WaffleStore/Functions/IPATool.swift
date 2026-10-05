//
//  IPATool.swift
//  PancakeStore
//
//  Created by Mineek on 19/10/2024.
//

// Heavily inspired by ipatool-py.
// https://github.com/NyaMisty/ipatool-py

import Foundation
import Zip
import SwiftUI
import PartyUI
import MapleSAP

typealias DownloadProgressHandler = (_ progress: Double, _ detail: String) -> Void

extension Data {
    var hexString: String {
        return map { String(format: "%02x", $0) }.joined()
    }
}

extension String {
    subscript (i: Int) -> String {
        return String(self[index(startIndex, offsetBy: i)])
    }

    subscript (r: Range<Int>) -> String {
        let start = index(startIndex, offsetBy: r.lowerBound)
        let end = index(startIndex, offsetBy: r.upperBound)
        return String(self[start..<end])
    }
}

class StoreClient {
    var session: URLSession
    var appleId: String
    var guid: String?
    var accountName: String?
    var authHeaders: [String: String]?
    var authCookies: [HTTPCookie]?
    var pod: String?

    init(account: StoreAccount) {
        let configuration = URLSessionConfiguration.ephemeral
        session = URLSession(configuration: configuration)
        appleId = account.email
        guid = account.guid
        accountName = account.name
        authHeaders = ["X-Dsid": account.dsid, "iCloud-Dsid": account.dsid,
            "X-Apple-Store-Front": account.storefront, "X-Token": account.passwordToken]
        authCookies = account.cookies.compactMap { $0.cookie() }
        pod = account.pod
    }

    func close() {
        session.invalidateAndCancel()
        authHeaders = nil; authCookies = nil
        guid = nil; pod = nil; accountName = nil; appleId = ""
    }

    func volumeStoreDownloadProduct(appId: String, appVerId: String = "") -> [String: Any] {
        var req = [
            "creditDisplay": "",
            "guid": self.guid!,
            "salableAdamId": appId,
        ]
        if appVerId != "" {
            req["externalVersionId"] = appVerId
        }
        let url = URL(string: "https://p\(pod!)-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/volumeStoreDownloadProduct?guid=\(self.guid!)")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.allHTTPHeaderFields = [
            "Content-Type": "application/x-www-form-urlencoded",
            "User-Agent": "Configurator/2.17 (Macintosh; OS X 15.2; 24C5089c) AppleWebKit/0620.1.16.11.6"
        ]
        let bodyString = req.map { "\($0.key)=\($0.value.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? "")" }.joined(separator: "&")
        request.httpBody = bodyString.data(using: .utf8)
        print("Setting headers")
        for (key, value) in self.authHeaders! {
            print("Setting authenticated header: \(key) [redacted]")
            request.addValue(value, forHTTPHeaderField: key)
        }
        print("Setting cookies")
        self.session.configuration.httpCookieStorage?.setCookies(self.authCookies!, for: url, mainDocumentURL: nil)

        var resp = [String: Any]()
        let datatask = session.dataTask(with: request) { (data, response, error) in
            if let error = error {
                print("error 2 \(error.localizedDescription)")
                return
            }
            if let data = data {
                do {
                    print("Got response")
                    let resp1 = try PropertyListSerialization.propertyList(from: data, options: [], format: nil) as! [String: Any]
                    if resp1["cancel-purchase-batch"] != nil {
                        print("Failed to download product: \(resp1["customerMessage"] as! String)")
                    }
                    resp = resp1
                } catch {
                    print("Error: \(error)")
                }
            }
        }
        datatask.resume()
        while datatask.state != .completed {
            sleep(1)
        }
        print("Got download response")
        return resp
    }

    func download(appId: String, appVer: String = "", isRedownload: Bool = false) -> [String: Any] {
        return self.volumeStoreDownloadProduct(appId: appId, appVerId: appVer)
    }

    func downloadToPath(url: String, path: String, progressHandler: DownloadProgressHandler? = nil) -> Void {
        var req = URLRequest(url: URL(string: url)!)
        req.httpMethod = "GET"
        let datatask = session.downloadTask(with: req) { (temporaryURL, response, error) in
            if let error = error {
                print("error 3 \(error.localizedDescription)")
                return
            }
            if let temporaryURL = temporaryURL {
                do {
                    let destinationURL = URL(fileURLWithPath: path)
                    if FileManager.default.fileExists(atPath: destinationURL.path) {
                        try FileManager.default.removeItem(at: destinationURL)
                    }
                    try FileManager.default.moveItem(at: temporaryURL, to: destinationURL)
                } catch {
                    print("Error: \(error)")
                }
            }
        }
        datatask.resume()
        while datatask.state != .completed {
            let progress = datatask.progress
            if progress.totalUnitCount > 0 {
                let fraction = min(max(progress.fractionCompleted, 0), 1)
                progressHandler?(fraction, String(format: "Downloading IPA %@".localized, "\(Int(fraction * 100))%"))
            }
            sleep(1)
        }
        progressHandler?(1, "Download complete".localized)
        print("Downloaded to \(path)")
    }
}

class IPATool {
    let storeClient: StoreClient

    init(account: StoreAccount) {
        storeClient = StoreClient(account: account)
    }

    func getVersionIDList(appId: String) -> [String] {
        print("Retrieving download info for appId \(appId)...")
        let downResp = storeClient.download(appId: appId, isRedownload: true)
        let songList = downResp["songList"] as? [[String: Any]] ?? []
        if songList.count == 0 {
            print("Failed to get id list!")
            return []
        }
        let downInfo = songList[0]
        let metadata = downInfo["metadata"] as? [String: Any] ?? [:]
        let appVerIds = metadata["softwareVersionExternalIdentifiers"] as? [Int] ?? []
        print("Got available version ids: \(appVerIds)")
        return appVerIds.map { String($0) }
    }

    func downloadIPAForVersion(appId: String, appVerId: String, progressHandler: DownloadProgressHandler? = nil) -> String {
        print("Downloading IPA for app \(appId) version \(appVerId)")
        progressHandler?(0.05, "Requesting download info".localized)
        let downResp = storeClient.download(appId: appId, appVer: appVerId)
        let songList = downResp["songList"] as! [[String: Any]]
        if songList.count == 0 {
            print("Failed to get app download info!")
            return ""
        }
        let downInfo = songList[0]
        let url = downInfo["URL"] as! String
        print("Received download URL [redacted]")
        let fm = FileManager.default
        let tempDir = fm.temporaryDirectory
        let path = tempDir.appendingPathComponent("app.ipa").path
        if fm.fileExists(atPath: path) {
            print("Removing existing file at \(path)")
            try! fm.removeItem(atPath: path)
        }
        storeClient.downloadToPath(url: url, path: path) { progress, detail in
            progressHandler?(0.10 + (progress * 0.60), detail)
        }
        Zip.addCustomFileExtension("ipa")
        progressHandler?(0.72, "Extracting IPA".localized)
        sleep(3)
        let path3 = URL(string: path)!
        let fileExtension = path3.pathExtension
        let fileName = path3.lastPathComponent
        let directoryName = fileName.replacingOccurrences(of: ".\(fileExtension)", with: "")
        let documentsUrl = fm.urls(for: .documentDirectory, in: .userDomainMask)[0]
        let destinationUrl = documentsUrl.appendingPathComponent(directoryName, isDirectory: true)
        if fm.fileExists(atPath: destinationUrl.path) {
            print("Removing existing folder at \(destinationUrl.path)")
            try! fm.removeItem(at: destinationUrl)
        }
        
        let unzipDirectory = try! Zip.quickUnzipFile(URL(string: path)!)
        progressHandler?(0.80, "Writing metadata".localized)
        var metadata = downInfo["metadata"] as! [String: Any]
        let metadataPath = unzipDirectory.appendingPathComponent("iTunesMetadata.plist").path
        metadata["apple-id"] = appleId
        metadata["userName"] = appleId
        (metadata as NSDictionary).write(toFile: metadataPath, atomically: true)
        print("Wrote iTunesMetadata.plist")
        var appContentDir = ""
        let payloadDir = unzipDirectory.appendingPathComponent("Payload")
        for entry in try! fm.contentsOfDirectory(atPath: payloadDir.path) {
            if entry.hasSuffix(".app") {
                print("Found app content dir: \(entry)")
                appContentDir = "Payload/" + entry
                break
            }
        }
        print("Found app content dir: \(appContentDir)")
        let scManifestData = try! Data(contentsOf: unzipDirectory.appendingPathComponent(appContentDir).appendingPathComponent("SC_Info").appendingPathComponent("Manifest.plist"))
        let scManifest = try! PropertyListSerialization.propertyList(from: scManifestData, options: [], format: nil) as! [String: Any]
        let sinfsDict = downInfo["sinfs"] as! [[String: Any]]
        if let sinfPaths = scManifest["SinfPaths"] as? [String] {
            progressHandler?(0.86, "Applying purchase data".localized)
            for (i, sinfPath) in sinfPaths.enumerated() {
                let sinfData = sinfsDict[i]["sinf"] as! Data
                try! sinfData.write(to: unzipDirectory.appendingPathComponent(appContentDir).appendingPathComponent(sinfPath))
                print("Wrote sinf to \(sinfPath)")
            }
        } else {
            print("Manifest.plist does not exist! Assuming it is an old app without one...")
            progressHandler?(0.86, "Applying purchase data".localized)
            let infoListData = try! Data(contentsOf: unzipDirectory.appendingPathComponent(appContentDir).appendingPathComponent("Info.plist"))
            let infoList = try! PropertyListSerialization.propertyList(from: infoListData, options: [], format: nil) as! [String: Any]
            let sinfPath = appContentDir + "/SC_Info/" + (infoList["CFBundleExecutable"] as! String) + ".sinf"
            let sinfData = sinfsDict[0]["sinf"] as! Data
            try! sinfData.write(to: unzipDirectory.appendingPathComponent(sinfPath))
            print("Wrote sinf to \(sinfPath)")
        }
        print("Downloaded IPA to \(unzipDirectory.path)")
        progressHandler?(0.90, "IPA prepared".localized)
        return unzipDirectory.path
    }
}

