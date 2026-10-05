import SwiftUI
import SafariServices
import Telegraph
import MapleSAP

// Original WaffleStore OTA mechanism, confined to loopback and one verified IPA.
// Opening Safari, or serving bytes, never proves that iOS installed the app.
@MainActor
final class OTAInstaller: ObservableObject {
    @Published var pageURL: URL?
    @Published var status = ""
    private var server: Server?
    private var directory: URL?
    private var expiration: Task<Void, Never>?
    private var backgroundTask: UIBackgroundTaskIdentifier = .invalid
    func start(_ record: DownloadRecord) throws {
        stop()
        guard let source = record.fileURL else { throw CocoaError(.fileNoSuchFile) }
        let fm = FileManager.default
        let root = fm.temporaryDirectory.appendingPathComponent("ota-" + UUID().uuidString, isDirectory: true)
        try fm.createDirectory(at: root, withIntermediateDirectories: true)
        directory = root
        do {
            let ipa = root.appendingPathComponent("signed.ipa")
            do { try fm.linkItem(at: source, to: ipa) } catch { try fm.copyItem(at: source, to: ipa) }
            let token = UUID().uuidString
            let path = "/" + token
            let fetch = "http://127.0.0.1:9090\(path)/ipa/signed.ipa"
            // Same HTTPS generator as upstream. Sends app metadata + loopback
            // address only; never Apple account/session data or uploaded IPA bytes.
            var manifest = URLComponents(string: "https://api.palera.in/genPlist")!
            manifest.queryItems = [URLQueryItem(name: "bundleid", value: record.bundleID),
                URLQueryItem(name: "name", value: record.appName),
                URLQueryItem(name: "version", value: record.build ?? record.version),
                URLQueryItem(name: "fetchurl", value: fetch)]
            var install = URLComponents(string: "itms-services://")!
            install.queryItems = [URLQueryItem(name: "action", value: "download-manifest"),
                URLQueryItem(name: "url", value: manifest.url!.absoluteString)]
            let target = install.url!.absoluteString
            // JSON escaping prevents metadata from becoming HTML/script content.
            let literal = String(data: try JSONSerialization.data(withJSONObject: target, options: .fragmentsAllowed), encoding: .utf8)!
            let html = "<html><meta name='viewport' content='width=device-width'><body><h3>Install selected version</h3><p>Confirm the iOS installation prompt. This page does not report installation success.</p><button id='install'>Install</button><script>const target=\(literal);document.getElementById('install').onclick=()=>location.href=target;location.href=target;</script></body></html>"
            let server = Server()
            server.concurrency = 2
            let fileHandler: HTTPRequest.Handler = { request in
                // Read-only mapping avoids allocating a full multi-GB IPA body.
                let mapped = try Data(contentsOf: ipa, options: .alwaysMapped)
                let range: Range<Int>
                do { range = try OTAByteRange.resolve(request.headers.range, size: mapped.count) }
                catch { return HTTPResponse(.rangeNotSatisfiable, headers: ["Content-Range": "bytes */\(mapped.count)"]) }
                let response = HTTPResponse(request.headers.range == nil ? .ok : .partialContent,
                    headers: ["Content-Type": "application/octet-stream", "Accept-Ranges": "bytes", "Cache-Control": "no-store"],
                    body: mapped[range])
                if request.headers.range != nil {
                    response.headers.contentRange = "bytes \(range.lowerBound)-\(range.upperBound - 1)/\(mapped.count)"
                }
                return response
            }
            server.route(.GET, path + "/ipa/signed.ipa", fileHandler)
            server.route(.HEAD, path + "/ipa/signed.ipa", fileHandler)
            server.route(.GET, path + "/install") { _ in HTTPResponse(.ok, headers: ["Content-Type": "text/html; charset=utf-8", "Cache-Control": "no-store"], content: html) }
            try server.start(port: 9090, interface: "127.0.0.1")
            self.server = server
            backgroundTask = UIApplication.shared.beginBackgroundTask(withName: "WaffleStore OTA transfer") {
                Task { @MainActor in
                    self.stop(); self.status = "iOS ended the background transfer time. The IPA was retained."
                }
            }
            pageURL = URL(string: "http://127.0.0.1:9090\(path)/install")!
            status = "Serving the verified IPA. iOS decides whether installation is allowed."
            expiration = Task {
                try? await Task.sleep(for: .seconds(600))
                guard !Task.isCancelled else { return }
                stop(); status = "Installation server timed out. The downloaded IPA was retained."
            }
        } catch { stop(); throw error }
    }
    func stop() {
        expiration?.cancel(); expiration = nil
        if backgroundTask != .invalid { UIApplication.shared.endBackgroundTask(backgroundTask); backgroundTask = .invalid }
        server?.stop(immediately: true); server = nil
        if let directory { try? FileManager.default.removeItem(at: directory) }
        directory = nil; pageURL = nil
    }
}
struct OTAInstallationView: View {
    let record: DownloadRecord
    @StateObject private var installer = OTAInstaller()
    @Environment(\.dismiss) private var dismiss
    @State private var started = false
    var body: some View {
        NavigationStack {
            Group {
                if let url = installer.pageURL { InstallationSafari(url: url) }
                else {
                    Form {
                        Section("Verified download") {
                            Text(record.appName)
                            Text("Version \(record.version) · ID \(record.externalVersionID)")
                            if let build = record.build { Text("IPA build \(build)") }
                            Text(record.bundleID).font(.caption)
                        }
                        Text("Uses WaffleStore's original Safari installation method and api.palera.in for the HTTPS manifest. It sends the app name, bundle ID, build and loopback URL; your account and IPA are not uploaded. Keep this screen open. iOS may refuse protected App Store packages or downgrades.")
                            .font(.caption)
                        if !installer.status.isEmpty { Text(installer.status) }
                        Button("Request iOS installation") {
                            do { try installer.start(record); started = true }
                            catch { installer.status = "Could not start local installation server (code \((error as NSError).code))." }
                        }
                        if let url = record.fileURL { ShareLink(item: url) { Label("Export IPA instead", systemImage: "square.and.arrow.up") } }
                    }
                }
            }
            .navigationTitle(started ? "Installation requested" : "Install downloaded IPA")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { installer.stop(); dismiss() } } }
            .onDisappear { installer.stop() }
        }
    }
}
private struct InstallationSafari: UIViewControllerRepresentable {
    let url: URL
    func makeUIViewController(context: Context) -> SFSafariViewController { SFSafariViewController(url: url) }
    func updateUIViewController(_ uiViewController: SFSafariViewController, context: Context) {}
}
