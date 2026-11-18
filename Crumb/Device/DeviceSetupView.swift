import SwiftUI

struct DeviceSetupView: View {
    var body: some View {
        Form {
            Section("On-device decryption") {
                Text("Decrypt downloaded apps directly in Crumb. Keep Crumb open until verification finishes.")
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
