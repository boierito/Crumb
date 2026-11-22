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

