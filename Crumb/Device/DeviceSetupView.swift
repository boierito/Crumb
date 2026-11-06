import SwiftUI
import UniformTypeIdentifiers

struct DeviceSetupView: View {
    @StateObject private var device = DeviceWorkflow.shared
    @State private var importKind = "p12"
    @State private var importing = false
    @Environment(\.scenePhase) private var scenePhase

    var body: some View {
        Form {
            Section {
                HStack(spacing: 14) {
                    Image("CrumbMark").resizable().frame(width: 52, height: 52).clipShape(RoundedRectangle(cornerRadius: 12))
                    VStack(alignment: .leading, spacing: 3) {
                        Text("Crumb Decrypter").font(.headline)
                        Text((device.ready && device.connected ? "Ready" : "Set up once").localized).font(.subheadline).foregroundStyle(.secondary)
                    }
                    Spacer()
                    if device.ready && device.connected { Image(systemName:"checkmark.circle.fill").foregroundStyle(.green) }
                }.padding(.vertical, 4)
            }
            Section("Connection") {
                HStack {
                    Label("LocalDevVPN", systemImage:"network")
                    Spacer()
                    if device.paired && device.busy && !device.signingReady { ProgressView() }
                    Text((device.connected ? "Connected" : "Waiting").localized).foregroundStyle(device.connected ? .green : .secondary)
                }
                HStack {
                    Label("Pairing", systemImage:"link")
                    Spacer()
                    if device.paired {
                        Image(systemName:"checkmark.circle.fill").foregroundStyle(.green)
                    } else if DeviceWorkflow.pairingSupported {
                        Button((device.busy ? "Cancel" : "Pair").localized) { if device.busy { device.cancel() } else { device.pair() } }
                    } else {
                        Button("Import") { pick("pairing") }.disabled(device.busy || device.waitingForHelper)
                    }
                }
                if DeviceWorkflow.pairingSupported && (!device.paired || device.pin != nil) {
                    Text("Settings → Privacy & Security → Developer Mode → Crumb").font(.footnote).foregroundStyle(.secondary)
                }
                if let pin = device.pin {
                    Text(pin).font(.system(.largeTitle, design:.rounded).monospacedDigit()).frame(maxWidth:.infinity).textSelection(.enabled)
                }
                if device.onCellular {
                    Label("Connect to Wi-Fi to pair this iPhone.", systemImage:"wifi.exclamationmark").font(.footnote).foregroundStyle(.orange)
                }
                if device.paired {
                    Menu("Change pairing") {
                        if DeviceWorkflow.pairingSupported { Button("Pair again") { device.pair() } }
                        Button("Import pairing file") { pick("pairing") }
                        Button("Remove pairing", role:.destructive) { device.forgetPairing() }
                    }.disabled(device.busy || device.waitingForHelper)
                }
            }
            SigningIdentitySections()
        }
        .navigationTitle("Decrypter")
        .navigationBarTitleDisplayMode(.inline)
        .onAppear { device.setupVisible = true }
        .onDisappear { device.setupVisible = false }
        .alert("Could not finish",isPresented:Binding(get:{ !device.errorMessage.isEmpty },set:{ if !$0 { device.errorMessage = "" } })) {
            Button("OK") { device.errorMessage = "" }
        } message: { Text(device.errorMessage) }
        .onChange(of:scenePhase) { phase in if phase == .active { device.autoConnect() } }
        .task { device.autoConnect() }
        .fileImporter(isPresented:$importing,allowedContentTypes:[.data],allowsMultipleSelection:false) { result in
            if case .success(let urls) = result, let url = urls.first { device.importFile(url,kind:importKind) }
        }
    }
    private func pick(_ kind: String) { importKind = kind; importing = true }

}

struct DeviceDecryptView: View {
    let record: DownloadRecord
    @StateObject private var device = DeviceWorkflow.shared
    var body: some View {
        Form {
            Section { Text(record.appName).font(.headline); Text(record.version).foregroundStyle(.secondary) }
            Section {
                NavigationLink("Decrypter settings") { DeviceSetupView() }
                Button("Decrypt") { device.decrypt(record) }.disabled(device.busy || !DeviceWorkflow.decryptSupported)
                Button("Decrypt and install") { device.decrypt(record,install:true) }.disabled(device.busy || !DeviceWorkflow.decryptSupported)
            }
        }.navigationTitle("Decrypt app").navigationBarTitleDisplayMode(.inline)
    }
}
