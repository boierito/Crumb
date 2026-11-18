import SwiftUI
import MapleSAP

struct ITunesSearchResponse: Decodable {
    let results: [ITunesApp]
}

struct ITunesApp: Decodable, Hashable, Identifiable {
    var id: Int { trackId }
    let trackId: Int
    let trackName: String
    let trackViewUrl: String?
    let artworkUrl100: String?
    let artistName: String?
    let bundleId: String?
    let version: String?
    let screenshotUrls: [String]?
    let ipadScreenshotUrls: [String]?
}

struct CatalogListing: Hashable {
    let input: String
    let country: String
    var app: ITunesApp? = nil
    var favourite: FavouriteApp? = nil

    static func request(_ endpoint: String, country: String, items: [URLQueryItem]) async throws -> [ITunesApp] {
        var components = URLComponents()
        components.scheme = "https"
        components.host = "itunes.apple.com"
        components.path = "/" + endpoint
        components.queryItems = items + [URLQueryItem(name: "country", value: country)]
        guard let url = components.url else { throw URLError(.badURL) }
        let (data, response) = try await URLSession.shared.data(from: url)
        guard let http = response as? HTTPURLResponse, (200..<300).contains(http.statusCode) else {
            throw URLError(.badServerResponse)
        }
        return try JSONDecoder().decode(ITunesSearchResponse.self, from: data).results
    }
}


struct AppSearchView: View {
    @EnvironmentObject private var appData: AppData
    @AppStorage("crumbSearchHistory") private var historyData = Data()
    @State private var query = ""
    @State private var results: [ITunesApp] = []
    @State private var isSearching = false
    @State private var searchError = false
    @State private var showRegion = false
    @State private var path: [CatalogListing] = []
    @State private var retry = 0
    @FocusState private var searchFocused: Bool

    private var term: String { query.trimmingCharacters(in: .whitespacesAndNewlines) }
    private var history: [String] { (try? JSONDecoder().decode([String].self, from: historyData)) ?? [] }
    private var directInput: Bool {
        !extractAppId(from: term).isEmpty ||
            (term.contains(".") && !term.contains(" ") && !term.contains(":") && term.split(separator: ".").count >= 3)
    }

    var body: some View {
        NavigationStack(path: $path) {
            List {
                Section {
                    VStack(alignment: .leading, spacing: 12) {
                        HStack(spacing: 10) {
                            Image(systemName: "magnifyingglass").foregroundStyle(.secondary).accessibilityHidden(true)
                            TextField("App name, link, ID or bundle ID", text: $query)
                                .focused($searchFocused)
                                .autocorrectionDisabled().textInputAutocapitalization(.never)
                                .submitLabel(directInput ? .go : .search)
                                .onSubmit {
                                    if directInput { open(CatalogListing(input: term, country: appData.catalogCountry)) }
                                    else { rememberTerm(); searchFocused = false }
                                }
                            if !query.isEmpty {
                                Button { query = "" } label: {
                                    Image(systemName: "xmark.circle.fill").frame(minWidth: 44, minHeight: 44)
                                }.buttonStyle(.borderless).accessibilityLabel("Clear search")
                            }
                        }
                        .padding(.horizontal, 12)
                        .frame(minHeight: 48)
                        .background(Color(uiColor: .tertiarySystemFill), in: RoundedRectangle(cornerRadius: 15))
                        Button { showRegion = true } label: {
                            Label(regionName(appData.catalogCountry), systemImage: "globe")
                                .font(.subheadline.weight(.medium))
                                .padding(.horizontal, 11).frame(minHeight: 44)
                                .background(Color.accentColor.opacity(0.12), in: Capsule())
                        }
                        .buttonStyle(.borderless)
                        .accessibilityLabel("Catalog region")
                        .accessibilityValue(regionName(appData.catalogCountry))
                    }
                    .padding(.vertical, 6)
                    .id(AppTab.search)
                }
                if directInput {
                    Section {
                        Button { open(CatalogListing(input: term, country: appData.catalogCountry)) } label: {
                            Label("Open app", systemImage: "arrow.up.right")
                        }.buttonStyle(CrumbActionStyle())
                    }
                } else if isSearching {
                    Section {
                        Label("Searching App Store…", systemImage: "magnifyingglass")
                            .foregroundStyle(.secondary)
                        ForEach(0..<2) { _ in
                            HStack {
                                CrumbArtwork(symbol: "app")
                                RoundedRectangle(cornerRadius: 4).fill(Color.secondary.opacity(0.15))
                                    .frame(width: 110, height: 10)
                            }.accessibilityHidden(true)
                        }
                    }
                } else if searchError {
                    Section {
                        CrumbEmptyState(title: "Search unavailable", message: "Check your connection and try again.", icon: "wifi.exclamationmark")
                        Button("Try again") { retry += 1 }
                    }
                } else if results.isEmpty {
                    if term.isEmpty && !history.isEmpty {
                        Section {
                            ForEach(history, id: \.self) { item in
                                Button { query = item; searchFocused = false } label: {
                                    Label(item, systemImage: "clock")
                                }
                            }
                            Button("Clear recent searches", role: .destructive) { historyData = Data() }
                        } header: { Text("Recent searches") }
                    } else {
                        CrumbEmptyState(title: term.isEmpty ? "Find an app" : "No apps found",
                                        message: term.isEmpty ? "Search the App Store or paste a link." : "Try another name, App ID or region.",
                                        icon: "magnifyingglass")
                    }
                } else {
                    Section {
                        ForEach(results) { app in
                            Button {
                                open(CatalogListing(input: app.trackViewUrl ?? String(app.trackId), country: appData.catalogCountry, app: app))
                            } label: {
                                HStack(spacing: 12) {
                                    CrumbArtwork(url: app.artworkUrl100, symbol: "app")
                                    VStack(alignment: .leading, spacing: 4) {
                                        Text(app.trackName).font(.headline).foregroundStyle(.primary)
                                        if let artist = app.artistName {
                                            Text(artist).font(.subheadline).foregroundStyle(.secondary)
                                        }
                                    }
                                    Spacer(minLength: 4)
                                    Image(systemName: "chevron.right").font(.caption.weight(.semibold)).foregroundStyle(.tertiary)
                                        .accessibilityHidden(true)
                                }.padding(.vertical, 4).contentShape(Rectangle())
                            }.buttonStyle(.plain)
                        }
                    } header: { Text("\(results.count) results") }
                }
            }
            .crumbPrimaryPage(.search)
            .navigationTitle("Search")
            .scrollDismissesKeyboard(.interactively)
            .sheet(isPresented: $showRegion) { CatalogRegionView() }
            .navigationDestination(for: CatalogListing.self) { CatalogDetailView(listing: $0) }
            .onChange(of: appData.linkedListing) { listing in
                if let listing { path = [listing]; appData.linkedListing = nil }
            }
            .onAppear {
                if let listing = appData.linkedListing { path = [listing]; appData.linkedListing = nil }
            }
            .task(id: "\(appData.catalogCountry)|\(retry)|\(term)") { await search() }
        }
    }

    private func open(_ listing: CatalogListing) {
        rememberTerm()
        searchFocused = false
        path.append(listing)
    }

    private func rememberTerm() {
        guard !term.isEmpty, !directInput else { return }
        historyData = (try? JSONEncoder().encode(Array(([term] + history.filter { $0 != term }).prefix(8)))) ?? Data()
    }

    private func search() async {
        results = []
        searchError = false
        isSearching = false
        guard !term.isEmpty, !directInput else { return }
        isSearching = true
        do {
            try await Task.sleep(nanoseconds: 300_000_000)
            let matches = try await CatalogListing.request("search", country: appData.catalogCountry, items: [
                URLQueryItem(name: "term", value: term),
                URLQueryItem(name: "entity", value: "software"),
                URLQueryItem(name: "limit", value: "25")
            ])
            try Task.checkCancellation()
            results = matches
        } catch {
            guard !Task.isCancelled else { return }
            searchError = true
        }
        isSearching = false
    }
}

struct CatalogDetailView: View {
    let listing: CatalogListing
    @EnvironmentObject private var appData: AppData
    @Environment(\.dynamicTypeSize) private var textSize
    @State private var resolved: ITunesApp?

    private var app: ITunesApp? { listing.app ?? resolved }
    private var name: String { app?.trackName ?? listing.favourite?.appName ?? "App Store app" }
    private var bundleID: String? { app?.bundleId ?? listing.favourite?.bundleId }
    private var saved: FavouriteApp? { appData.favourites.first { $0.bundleId == bundleID } }
    private var artwork: String? { app?.artworkUrl100 ?? listing.favourite?.artworkURL }
    private var previews: [String] { app?.screenshotUrls ?? [] }

    var body: some View {
        List {
            Section {
                let layout = textSize.isAccessibilitySize ? AnyLayout(VStackLayout(alignment: .leading, spacing: 12)) : AnyLayout(HStackLayout(spacing: 16))
                layout {
                    CrumbArtwork(url: artwork, symbol: "app", size: 78)
                    VStack(alignment: .leading, spacing: 5) {
                        Text(name).font(.title2.weight(.bold)).accessibilityAddTraits(.isHeader)
                        if let artist = app?.artistName { Text(artist).foregroundStyle(.secondary) }
                        Text(regionName(listing.country)).font(.footnote).foregroundStyle(.secondary)
                    }
                    if let bundleID {
                        Button {
                            if let saved { appData.favourites = FavouritesStore.remove(saved) }
                            else {
                                appData.favourites = FavouritesStore.add(FavouriteApp(appLink: listing.input,
                                    bundleId: bundleID, appName: name, catalogCountry: listing.country, artworkURL: artwork))
                            }
                        } label: {
                            Image(systemName: saved == nil ? "star" : "star.fill").frame(width: 44, height: 44)
                        }
                        .buttonStyle(.borderless)
                        .accessibilityLabel(saved == nil ? "Add to favourites" : "Remove from favourites")
                    }
                }.padding(.vertical, 8)
                Button { open(latest: true) } label: {
                    Label("Download latest version", systemImage: "arrow.down.circle")
                }.buttonStyle(CrumbActionStyle())
                Button { open(latest: false) } label: {
                    Label("View versions", systemImage: "clock")
                }.buttonStyle(CrumbActionStyle(primary: false))
            } footer: {
                if appData.isAuthenticated, listing.country != appData.accountCountry {
                    Text("Apple checks download access against your account region (\(regionName(appData.accountCountry))).")
                }
            }
            if !previews.isEmpty {
                Section("Preview") {
                    ScrollView(.horizontal) {
                        HStack(spacing: 10) {
                            ForEach(Array(previews.prefix(5)), id: \.self) { url in
                                AsyncImage(url: URL(string: url)) { image in
                                    image.resizable().scaledToFit()
                                } placeholder: {
                                    RoundedRectangle(cornerRadius: 14).fill(Color(uiColor: .tertiarySystemFill))
                                }
                                .frame(width: 140, height: 240)
                                .clipShape(RoundedRectangle(cornerRadius: 14))
                                .accessibilityHidden(true)
                            }
                        }
                    }
                }
            }
        }
        .navigationTitle("App details")
        .navigationBarTitleDisplayMode(.inline)
        .task {
            guard app == nil else { return }
            let id = extractAppId(from: listing.input)
            let item = id.isEmpty ? URLQueryItem(name: "bundleId", value: bundleID ?? listing.input)
                                  : URLQueryItem(name: "id", value: id)
            resolved = try? await CatalogListing.request("lookup", country: listing.country, items: [item]).first
        }
    }

    private func open(latest: Bool) {
        appData.openAppSelection(listing.input, country: listing.country, latest: latest)
    }
}

func regionName(_ code: String) -> String {
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
