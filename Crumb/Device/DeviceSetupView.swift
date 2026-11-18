import SwiftUI

struct DeviceSetupView: View {
    var body: some View {
        Form {
            Section {
                CrumbEmptyState(title: "Device tools",
                    message: "Manage package tools and the signing identity used for installation.",
                    icon: "iphone")
            }
            Section("Package tools") {
                Label(DeviceWorkflow.decryptSupported ? "Available on this device" : "Unavailable on this device",
                      systemImage: DeviceWorkflow.decryptSupported ? "checkmark.circle" : "info.circle")
                Text("Open a downloaded app’s options to start. Keep Crumb open until verification finishes.")
                    .font(.subheadline).foregroundStyle(.secondary)
            }
            Section("Installation") {
                NavigationLink { SigningSettingsView() } label: {
                    Label("Signing identity", systemImage: "checkmark.seal")
                }
            }
        }
        .navigationTitle("Device tools")
        .navigationBarTitleDisplayMode(.inline)
    }
}

struct DeviceDecryptView: View {
    let record: DownloadRecord
    @StateObject private var device = DeviceWorkflow.shared
    var body: some View {
        Form {
            Section {
                HStack(spacing: 12) {
                    CrumbArtwork()
                    VStack(alignment: .leading) {
                        Text(record.appName).font(.headline)
                        Text(record.version).foregroundStyle(.secondary)
                    }
                }
            }
            Section {
                NavigationLink("Decrypter settings") { DeviceSetupView() }
                Button("Decrypt") { device.decrypt(record) }.buttonStyle(CrumbActionStyle()).disabled(device.busy || !DeviceWorkflow.decryptSupported)
                Button("Decrypt and install") { device.decrypt(record,install:true) }.disabled(device.busy || !DeviceWorkflow.decryptSupported)
            }
        }.navigationTitle("Decrypt app").navigationBarTitleDisplayMode(.inline)
    }
}
