import SwiftUI

struct FavouritesView: View {
    @EnvironmentObject private var appData: AppData

    var body: some View {
        NavigationStack {
            List {
                if appData.favourites.isEmpty {
                    Section {
                        CrumbEmptyState(title: "No favourites yet", message: "Save an app from its details to find it here.", icon: "star")
                            .id(AppTab.favourites)
                        Button("Search App Store") { appData.selectedTab = .search }
                            .buttonStyle(CrumbActionStyle(primary: false))
                    }
                } else {
                    Text("\(appData.favourites.count) saved apps")
                        .font(.subheadline).foregroundStyle(.secondary).id(AppTab.favourites)
                    ForEach(appData.favourites) { favourite in
                        let country = favourite.catalogCountry ?? appData.accountCountry
                        NavigationLink(value: CatalogListing(input: favourite.appLink, country: country, favourite: favourite)) {
                            HStack(spacing: 12) {
                                CrumbArtwork(url: favourite.artworkURL, symbol: "app")
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(favourite.appName).font(.headline)
                                    Text(regionName(country)).font(.subheadline).foregroundStyle(.secondary)
                                }
                            }.padding(.vertical, 4)
                        }
                        .contextMenu { removeButton(favourite) }
                        .swipeActions(allowsFullSwipe: false) { removeButton(favourite) }
                    }
                }
            }
            .crumbPrimaryPage(.favourites)
            .navigationTitle("Favourites".localized)
            .navigationDestination(for: CatalogListing.self) { CatalogDetailView(listing: $0) }
        }
    }

    private func removeButton(_ favourite: FavouriteApp) -> some View {
        Button(role: .destructive) {
            appData.favourites = FavouritesStore.remove(favourite)
        } label: {
            Label("Remove favourite", systemImage: "star.slash")
        }
    }
}
