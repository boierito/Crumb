import SwiftUI
import UniformTypeIdentifiers

struct DeviceSetupView: View {
    @StateObject private var device = DeviceWorkflow.shared
    @State private var importKind = "p12"
    @State private var importing = false
    @Environment(\.scenePhase) private var scenePhase

    var body: some View {
        Form {
            Section("On-device decryption") {
                Label(AppGroupDecryption.container() == nil ? "App Group unavailable" : "App Group ready",
                      systemImage: AppGroupDecryption.container() == nil ? "exclamationmark.triangle" : "checkmark.circle")
                Text("Decrypt downloaded apps directly in Crumb. No helper, VPN or pairing is required. Keep Crumb open until verification finishes.")
                Text("In KSign, sign Crumb with an App Group authorized by your provisioning profile. The Apple account licensed for the app must be available to FairPlay on this device.")
                    .font(.footnote).foregroundStyle(.secondary)
                Text("Experimental: FairPlay may reject this method on your OS version. Unsupported FAT or 32-bit binaries fail without publishing a partial IPA.")
                    .font(.footnote).foregroundStyle(.secondary)
            }
            SigningIdentitySection()
        }
        .navigationTitle("Decrypter")
        .navigationBarTitleDisplayMode(.inline)
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
