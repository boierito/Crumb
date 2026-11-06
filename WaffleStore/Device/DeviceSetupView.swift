import SwiftUI
import UniformTypeIdentifiers

struct DeviceSetupView: View {
    @StateObject private var device = DeviceWorkflow.shared
    @State private var importKind = "p12"
    @State private var importing = false
    @State private var editingSigning = false
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
            Section("Signing") {
                if device.signingReady && !editingSigning {
                    Label(device.certificateName.isEmpty ? "Certificate" : device.certificateName, systemImage:"checkmark.seal.fill")
                        .lineLimit(1).foregroundStyle(.primary)
                    Label(device.profileName.isEmpty ? "Provisioning profile" : device.profileName, systemImage:"doc.badge.gearshape")
                        .lineLimit(1)
                    HStack { Label("Password saved", systemImage:"lock.fill"); Spacer(); Button("Change") { editingSigning = true } }
                        .font(.subheadline).foregroundStyle(.secondary)
                } else {
                    signingFile(title:"Certificate", name:device.certificateName, present:device.hasCertificate, kind:"p12")
                    signingFile(title:"Provisioning profile", name:device.profileName, present:device.hasProfile, kind:"profile")
                    SecureField("Certificate password", text:$device.signingPassword)
                        .textContentType(.password)
                        .onChange(of:device.signingPassword) { _ in device.invalidateSigning() }
                    Button { device.saveSigning() } label: {
                        HStack { Text("Save signing"); Spacer(); if device.busy && device.phase == "signing" { ProgressView() } }
                    }.disabled(device.busy || !device.hasCertificate || !device.hasProfile)
                }
            }
        }
        .navigationTitle("Decrypter")
        .navigationBarTitleDisplayMode(.inline)
        .onChange(of:device.signingReady) { ready in if ready { editingSigning = false } }
        .onChange(of:scenePhase) { phase in if phase == .active { device.autoConnect(force:true) } }
        .task {
            while !Task.isCancelled {
                device.autoConnect()
                try? await Task.sleep(for:.seconds(5))
            }
        }
        .fileImporter(isPresented:$importing,allowedContentTypes:[.data],allowsMultipleSelection:false) { result in
            if case .success(let urls) = result, let url = urls.first { device.importFile(url,kind:importKind) }
        }
    }
    private func pick(_ kind: String) { importKind = kind; importing = true }
    private func signingFile(title:String,name:String,present:Bool,kind:String) -> some View {
        Button { pick(kind) } label: {
            HStack(spacing:12) {
                Image(systemName:present ? "checkmark.circle.fill" : "plus.circle").foregroundStyle(present ? .green : .accentColor)
                VStack(alignment:.leading,spacing:2) {
                    Text(title.localized).foregroundStyle(.primary)
                    if !name.isEmpty { Text(name).font(.caption).foregroundStyle(.secondary).lineLimit(1) }
                }
                Spacer()
                Image(systemName:"chevron.right").font(.caption).foregroundStyle(.tertiary)
            }
        }.disabled(device.busy || device.waitingForHelper)
    }
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
