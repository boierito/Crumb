import SwiftUI
import PartyUI
import MapleSAP

struct ITunesSearchResponse: Codable {
    let resultCount: Int
    let results: [ITunesApp]
}
struct ITunesApp: Codable, Identifiable {
    var id: Int { trackId }
    let trackId: Int
    let trackName: String
    let trackViewUrl: String
    let artworkUrl100: String
    let artistName: String
    let primaryGenreName: String?
}

struct AppSearchView: View {
    var embedded = false
    @Environment(\.dismiss) private var dismiss
    @EnvironmentObject var appData: AppData
    @State private var query = ""
    @State private var results: [ITunesApp] = []
    @State private var isSearching = false
    @State private var searchError = ""
    @State private var searchTask: Task<Void, Never>?
    @State private var showRegion = false

    private var directInput: Bool {
        let value = query.trimmingCharacters(in: .whitespacesAndNewlines)
        return !extractAppId(from: value).isEmpty ||
            (value.contains(".") && !value.contains(" ") && !value.contains(":") && value.split(separator: ".").count >= 3)
    }
    var body: some View {
        NavigationStack {
            List {
                Section {
                    HStack {
                        Image(systemName: "magnifyingglass").foregroundStyle(.secondary)
                        TextField("App name, link, ID or bundle ID", text: $query)
                            .autocorrectionDisabled().textInputAutocapitalization(.never)
                            .submitLabel(directInput ? .go : .search)
                            .onSubmit { if directInput { choose(query) } }
                        if !query.isEmpty {
                            Button { query = "" } label: { Image(systemName: "xmark.circle.fill") }
                                .buttonStyle(.borderless).accessibilityLabel("Clear search")
                        }
                    }
                    if directInput {
                        Button { choose(query) } label: { Label("Choose version", systemImage: "arrow.down.app") }
                            .disabled(appData.isDowngrading || appData.storeRequestCount > 0)
                    }
                } footer: {
                    if appData.catalogCountry != appData.accountCountry {
                        Text("Browsing \(regionName(appData.catalogCountry)). Apple decides download access for your account (\(appData.accountCountry.uppercased())).")
                    }
                }
                if isSearching {
                    HStack(spacing: 12) { ProgressView(); Text("Searching App Store…").foregroundStyle(.secondary) }
                } else if !searchError.isEmpty {
                    Section {
                        Text(searchError).foregroundStyle(.secondary)
                        Button("Try again") { triggerSearch() }
                    }
                } else if results.isEmpty && !directInput {
                    VStack(alignment: .leading, spacing: 8) {
                        Label(query.isEmpty ? "Search App Store" : "No apps found", systemImage: "magnifyingglass")
                            .font(.headline)
                        Text(query.isEmpty ? "Find an app, then choose the version to download or install." : "Try another name or region.")
                            .font(.subheadline).foregroundStyle(.secondary)
                    }.padding(.vertical, 8)
                }
                ForEach(results) { app in
                    Button { choose(app.trackViewUrl) } label: {
                        HStack(spacing: 12) {
                            AsyncImage(url: URL(string: app.artworkUrl100)) { image in
                                image.resizable().aspectRatio(contentMode: .fit)
                            } placeholder: { Image(systemName: "app").foregroundStyle(.secondary) }
                            .frame(width: 50, height: 50).clipShape(RoundedRectangle(cornerRadius: 10))
                            VStack(alignment: .leading, spacing: 4) {
                                Text(app.trackName).font(.headline).foregroundStyle(.primary)
                                Text(app.artistName).font(.subheadline).foregroundStyle(.secondary)
                                if let genre = app.primaryGenreName { Text(genre).font(.caption).foregroundStyle(.secondary) }
                            }
                            Spacer(minLength: 0)
                            Image(systemName: "chevron.right").font(.caption).foregroundStyle(.tertiary)
                        }.padding(.vertical, 4)
                    }
                    .buttonStyle(.plain)
                    .disabled(appData.isDowngrading || appData.storeRequestCount > 0)
                }
                if appData.isDowngrading {
                    Section {
                        Button { appData.selectedTab = .downloads } label: {
                            Label("View download progress", systemImage: "square.and.arrow.down")
                        }
                    }
                }
            }
            .navigationTitle("Search")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                if !embedded { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
                ToolbarItem(placement: .topBarTrailing) {
                    Button { showRegion = true } label: {
                        Label(appData.catalogCountry.uppercased(), systemImage: "globe")
                    }
                    .accessibilityLabel("Catalog region: \(regionName(appData.catalogCountry))")
                }
            }
            .sheet(isPresented: $showRegion) { CatalogRegionView() }
            .onChange(of: query) { _ in triggerSearch() }
            .onChange(of: appData.catalogCountry) { _ in triggerSearch() }
            .onAppear { if results.isEmpty { triggerSearch() } }
            .onDisappear { searchTask?.cancel(); isSearching = false }
        }
    }
    private func choose(_ input: String) {
        Haptic.shared.play(.light)
        if embedded { appData.openAppSelection(input.trimmingCharacters(in: .whitespacesAndNewlines)) }
        else { appData.appLink = input; dismiss() }
    }
    private func triggerSearch() {
        searchTask?.cancel()
        results = []; searchError = ""; isSearching = false
        let term = query.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !term.isEmpty, !directInput else { return }
        let country = appData.catalogCountry
        isSearching = true
        searchTask = Task {
            do {
                try await Task.sleep(nanoseconds: 300_000_000)
                var components = URLComponents(string: "https://itunes.apple.com/search")!
                components.queryItems = [URLQueryItem(name: "term", value: term), URLQueryItem(name: "entity", value: "software"), URLQueryItem(name: "limit", value: "25"), URLQueryItem(name: "country", value: country)]
                guard let url = components.url else { throw URLError(.badURL) }
                let (data, response) = try await URLSession.shared.data(from: url)
                try Task.checkCancellation()
                guard let http = response as? HTTPURLResponse, (200..<300).contains(http.statusCode) else {
                    throw URLError(.badServerResponse)
                }
                results = try JSONDecoder().decode(ITunesSearchResponse.self, from: data).results
                isSearching = false
            } catch {
                guard !Task.isCancelled, (error as? URLError)?.code != .cancelled else { return }
                searchError = "Could not search App Store. Check your connection and try again."
                isSearching = false
            }
        }
    }
}

private func regionName(_ code: String) -> String {
    Locale.current.localizedString(forRegionCode: code.uppercased()) ?? code.uppercased()
}

private struct CatalogRegionView: View {
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @State private var query = ""
    private var countries: [String] {
        Storefront.catalogCountries.sorted { regionName($0).localizedStandardCompare(regionName($1)) == .orderedAscending }
            .filter { query.isEmpty || regionName($0).localizedCaseInsensitiveContains(query) || $0.localizedCaseInsensitiveContains(query) }
    }
    var body: some View {
        NavigationStack {
            List {
                Section {
                    row("", name: "Use account region (\(appData.accountCountry.uppercased()))")
                } footer: {
                    Text("Search and resolve versions in this catalog. You can try downloading with your current account; Apple decides license availability. This does not change your Apple account region.")
                }
                Section("Catalogs") {
                    ForEach(countries, id: \.self) { country in row(country, name: regionName(country)) }
                }
            }
            .searchable(text: $query, prompt: "Country or region")
            .navigationTitle("Catalog region").navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
        }
    }
    private func row(_ code: String, name: String) -> some View {
        Button {
            appData.catalogRegion = code
            dismiss()
        } label: {
            HStack {
                Text(name).foregroundStyle(.primary)
                Spacer()
                if code == appData.catalogRegion { Image(systemName: "checkmark") }
            }
        }
    }
}
