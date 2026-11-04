import SwiftUI
import Foundation
import UniformTypeIdentifiers

struct DeviceSetupView: View {
    @StateObject private var device = DeviceWorkflow.shared
    @State private var importKind = "p12"
    @State private var importing = false
    var body: some View {
        Form {
            Section {
                Text("Experimental device connection").font(.headline)
                Text("Pair from Settings → Privacy & Security → Developer Mode. Allow Local Network access. When only cellular data is available, use the no-Wi-Fi procedure below.")
                Text("This build enables pairing on iOS 27.0–27.0.1 and the decryption experiment on iOS 27.0.1. Physical validation in Crumb is pending.").font(.caption).foregroundStyle(.secondary)
            }
            Section("Pair this device") {
                if DeviceWorkflow.pairingSupported {
                    Button(device.paired ? "Pair again" : "Start pairing") { device.pair() }.disabled(device.busy)
                    if device.phase == "waiting" { Text("Open Developer Mode in Settings and select Crumb. Return here to see the PIN if needed.") }
                    if let pin = device.pin {
                        Text(pin).font(.largeTitle.monospacedDigit()).textSelection(.enabled)
                        Text("Enter this PIN in Settings to approve the pairing.")
                    }
                    Picker("Local device route", selection: $device.route) {
                        Text("LocalDevVPN").tag("10.7.0.1")
                        Text("Loopback (experimental)").tag("127.0.0.1")
                    }.disabled(device.busy)
                    if device.route == "127.0.0.1" {
                        Text("Direct access depends on iOS allowing this connection. If it fails, use LocalDevVPN.").font(.caption).foregroundStyle(.secondary)
                    }
                    Text("Pairing port: Bonjour discovery, saved port, or manual override. 49152 is a fallback.")
                        .font(.caption).foregroundStyle(.secondary)
                    Button("Discover and save RemotePairing port") { device.discoverPairingPort() }
                        .disabled(device.busy)
                    if let saved = device.lastPairingPort {
                        Text("Last discovered: \(saved) (may change after reboot).")
                            .font(.caption).foregroundStyle(.secondary)
                    }
                    TextField("Manual pairing port (optional)", text: $device.manualPairingPort)
                        .keyboardType(.numberPad)
                        .disabled(device.busy)
                    if !device.routeDiscovery.isEmpty {
                        Text(device.routeDiscovery).font(.caption).foregroundStyle(.secondary)
                    }
                    Button("Check authenticated connection") { device.verify() }.disabled(!device.paired || device.busy)
                    if device.connected { Label("Connected", systemImage: "checkmark.circle") }
                    Button("Remove local pairing keys", role: .destructive) { device.forgetPairing() }.disabled(device.busy || !device.paired)
                } else { Text("On-device pairing is unavailable on this system version in this build.") }
            }
            if DeviceWorkflow.pairingSupported {
                Section("No Wi-Fi? Use cellular data first") {
                    Text("1. Leave Cellular Data ON and connect LocalDevVPN.")
                    Text("2. With data still ON, tap Discover and save RemotePairing port above.")
                    Text("3. Open Control Center. Turn Cellular Data OFF, but leave LocalDevVPN connected.")
                    Text("4. Return to Crumb; tap Check authenticated connection, then Run VPN access diagnostics.")
                    Text("5. Turn Cellular Data ON again after testing.")
                    Text("iOS does not let Crumb switch cellular data directly. If the VPN disconnects, turn data on again before reconnecting it. This is an experimental workaround, not a guaranteed iOS permission bypass.")
                        .font(.caption).foregroundStyle(.secondary)
                }
                Section("No-helper VPN diagnostics") {
                    Text("Tests the existing authenticated LocalDevVPN connection, a harmless AFC canary in PublicStaging, Crumb's own Documents mapping, CoreDevice and a bounded passive OS trace sample. No IPA, SINF, helper or certificate is involved.")
                        .font(.caption).foregroundStyle(.secondary)
                    Button("Run VPN access diagnostics") { device.probeVPN() }
                        .disabled(device.busy || !device.connected)
                    if !device.diagnosticReport.isEmpty {
                        Text(device.diagnosticReport)
                            .font(.system(.caption, design: .monospaced))
                            .textSelection(.enabled)
                    }
                }
            }
            if DeviceWorkflow.decryptSupported {
                Section {
                    DisclosureGroup("Legacy helper signing (not used for VPN diagnostics)") {
                    Text("Import a .p12 certificate and a matching provisioning profile that permits your chosen helper identifier and this device. The certificate password is requested when preparing a download. Files are stored in this device's Keychain.").font(.caption)
                    Text("Crumb Decrypt is a separate app. A profile limited to Crumb's own identifier cannot sign this helper.").font(.caption).foregroundStyle(.secondary)
                    TextField("Helper bundle identifier", text: $device.helperBundle).textInputAutocapitalization(.never).autocorrectionDisabled().disabled(device.busy)
                    Button(device.hasCertificate ? "Replace certificate (.p12)" : "Import certificate (.p12)") { importKind = "p12"; importing = true }.disabled(device.busy)
                    Button(device.hasProfile ? "Replace provisioning profile" : "Import provisioning profile") { importKind = "profile"; importing = true }.disabled(device.busy)
                    Button("Remove signing files", role: .destructive) { device.forgetSigning() }.disabled(device.busy)
                    }
                }
            }
            DeviceStatusSection()
        }
        .navigationTitle("Device connection")
        .fileImporter(isPresented: $importing, allowedContentTypes: [.data], allowsMultipleSelection: false) { result in
            switch result {
            case .success(let urls): if let url = urls.first { device.importFile(url, kind: importKind) }
            case .failure: break
            }
        }
    }
}
struct DeviceStatusSection: View {
    @StateObject private var device = DeviceWorkflow.shared
    var body: some View {
        Section("Device operation") {
            if device.busy { ProgressView(device.phase.capitalized); Button("Cancel", role: .cancel) { device.cancel() } }
            if !device.message.isEmpty { Text(device.message).font(.callout) }
            if device.helperInstalled {
                Button("Open Crumb Decrypt") { device.openHelper() }.disabled(device.busy)
                Button("Retrieve result") { device.retrieve() }.disabled(device.busy)
            }
            if let output = device.output { ShareLink(item: output) { Label("Export decrypted IPA", systemImage: "square.and.arrow.up") } }
        }
    }
}
struct DeviceDecryptView: View {
    let record: DownloadRecord
    @StateObject private var device = DeviceWorkflow.shared
    @State private var password = ""
    var body: some View {
        Form {
            Section("Decrypt downloaded app") {
                Text(record.appName).font(.headline)
                Text("Version \(record.version)")
                Text("The device's App Store account must own this app. The helper preserves its Apple license and encrypted bytes until FairPlay processes them. Keep the helper in the foreground while decrypting.")
                Text("Experimental: iOS 27.0.1, thin arm64 packages, maximum 4 GB expanded and 512 MiB per file. A failure produces no completed IPA.").font(.caption).foregroundStyle(.secondary)
            }
            Section {
                NavigationLink("Connection and signing setup") { DeviceSetupView() }
                SecureField("Certificate password", text: $password)
                Button("Prepare and install decryption helper") {
                    device.prepare(record, password: password); password = ""
                }.disabled(device.busy || !device.connected || !device.hasCertificate || !device.hasProfile || !DeviceWorkflow.decryptSupported)
            }
            DeviceStatusSection()
        }
        .navigationTitle("On-device decryption")
        .onDisappear { password = "" }
    }
}
