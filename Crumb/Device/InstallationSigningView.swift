import SwiftUI
import UniformTypeIdentifiers

struct AdditionalSigningProfile: Codable, Identifiable {
    let id: String
    let name: String
    let data: Data
}
struct SigningReview: Identifiable {
    let id: UUID
    let appName: String
    let originalBundleID: String
    let bundleID: String
    let bundleCount: Int
    let changedIDs: Int
    let capabilities: [String]
    let profileName: String
}
struct SignedInstallation: Identifiable {
    let id = UUID()
    let record: DownloadRecord
    let source: URL
    let workspace: URL
}
struct InstallationSigningView: View {
    @StateObject private var device = DeviceWorkflow.shared
    @State private var importing = false
    @State private var kind = "install-p12"
    @State private var editing = false
    var body: some View {
        Form {
            Section("Signing identity") {
                if device.installationReady && !editing {
                    Label(device.installationCertificateName, systemImage:"checkmark.seal.fill").lineLimit(1)
                    Label(device.installationProfileName, systemImage:"doc.badge.gearshape").lineLimit(1)
                    HStack { Label("Password saved",systemImage:"lock.fill"); Spacer(); Button("Change") { editing = true } }
                } else {
                    fileRow("Certificate",name:device.installationCertificateName,kind:"install-p12")
                    fileRow("Provisioning profile",name:device.installationProfileName,kind:"install-profile")
                    SecureField("Certificate password",text:$device.installationPassword)
                        .onChange(of:device.installationPassword) { _ in device.invalidateInstallationSigning() }
                    Button("Save signing") { device.saveInstallationSigning() }
                        .disabled(device.busy || !device.hasInstallationCertificate || !device.hasInstallationProfile)
                }
            }
            Section("Extension profiles") {
                ForEach(device.installationProfiles) { profile in
                    HStack {
                        Text(profile.name).lineLimit(1)
                        Spacer()
                        Button(role:.destructive) { device.removeInstallationProfile(profile.id) } label: { Image(systemName:"minus.circle") }
                            .buttonStyle(.borderless).accessibilityLabel("Remove profile")
                    }
                }
                Button("Add profiles") { kind = "install-extra-profile"; importing = true }.disabled(device.busy)
            }
            if device.busy { Section { ProgressView("Checking identity…") } }
        }
        .navigationTitle("App installation").navigationBarTitleDisplayMode(.inline)
        .onAppear { device.setupVisible = true; device.seedInstallationSigning() }
        .onDisappear { device.setupVisible = false }
        .onChange(of:device.installationReady) { if $0 { editing = false } }
        .alert("Could not finish",isPresented:Binding(get:{ !device.errorMessage.isEmpty },set:{ if !$0 { device.errorMessage = "" } })) {
            Button("OK") { device.errorMessage = "" }
        } message: { Text(device.errorMessage) }
        .fileImporter(isPresented:$importing,allowedContentTypes:[.data],allowsMultipleSelection:kind == "install-extra-profile") { result in
            if case .success(let urls) = result { for url in urls { device.importFile(url,kind:kind) } }
        }
    }
    private func fileRow(_ title:String,name:String,kind:String)->some View {
        Button { self.kind = kind; importing = true } label: {
            HStack {
                Image(systemName:name.isEmpty ? "plus.circle" : "doc.fill")
                VStack(alignment:.leading) {
                    Text(title.localized).foregroundStyle(.primary)
                    if !name.isEmpty { Text(name).font(.caption).foregroundStyle(.secondary).lineLimit(1) }
                }
                Spacer(); Image(systemName:"chevron.right").foregroundStyle(.tertiary)
            }
        }.disabled(device.busy)
    }
}
struct SigningReviewView: View {
    let review: SigningReview
    @StateObject private var device = DeviceWorkflow.shared
    @State private var bundleID = ""
    var body: some View {
        NavigationStack {
            Form {
                Section { Text(review.appName).font(.headline); Text(review.profileName).foregroundStyle(.secondary) }
                Section("Bundle ID") {
                    TextField("Bundle ID",text:$bundleID).textInputAutocapitalization(.never).autocorrectionDisabled()
                    if review.originalBundleID != review.bundleID {
                        LabeledContent("Original",value:review.originalBundleID).font(.caption)
                    }
                }
                if review.bundleCount > 1 { Section { LabeledContent("App and extensions",value:String(review.bundleCount)) } }
                if !review.capabilities.isEmpty {
                    Section {
                        ForEach(review.capabilities,id:\.self) { key in Text(SigningSummary.capabilityName(key).localized) }
                    } header: { Text("Adjusted capabilities") }
                    footer: { Text("These features may be unavailable after signing.") }
                }
                Section {
                    Button((bundleID != review.bundleID ? "Review changes" : "Sign and install").localized) {
                        if bundleID != review.bundleID { device.replanSigning(bundleID:bundleID) }
                        else { device.approveSigning() }
                    }.disabled(device.busy || bundleID.isEmpty)
                    if device.busy { ProgressView("Checking signing…") }
                }
            }
            .disabled(device.busy)
            .navigationTitle("Sign app").navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement:.cancellationAction) { Button("Cancel") { device.stopDecryption() } } }
            .onAppear { bundleID = review.bundleID }
            .onChange(of:review.bundleID) { bundleID = $0 }
            .alert("Could not finish",isPresented:Binding(get:{ !device.errorMessage.isEmpty },set:{ if !$0 { device.errorMessage = "" } })) {
                Button("OK") { device.errorMessage = "" }
            } message: { Text(device.errorMessage) }
        }
    }
}
struct SigningDetailsView: View {
    let summary: SigningSummary
    var body: some View {
        Form {
            Section("Last signature") {
                LabeledContent("Bundle ID",value:summary.bundleID).font(.callout)
                Text(summary.profileName)
                Text(summary.date,style:.date)
            }
            if !summary.capabilities.isEmpty {
                Section("Adjusted capabilities") {
                    ForEach(summary.capabilities,id:\.self) { Text(SigningSummary.capabilityName($0).localized) }
                }
            }
            Section { ShareLink(item:summary.report) { Label("Share signing report",systemImage:"square.and.arrow.up") } }
        }.navigationTitle("Signing details")
    }
}
