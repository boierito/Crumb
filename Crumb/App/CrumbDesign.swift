import SwiftUI

struct CrumbArtwork: View {
    var url: String? = nil
    var symbol = "shippingbox"
    var size: CGFloat = 48

    var body: some View {
        AsyncImage(url: url.flatMap { URL(string: $0) }) { phase in
            if let image = phase.image {
                image.resizable().scaledToFill()
            } else {
                Image(systemName: symbol)
                    .font(.system(size: size * 0.40, weight: .medium))
                    .foregroundStyle(Color.accentColor)
                    .frame(width: size, height: size)
                    .background(Color.accentColor.opacity(0.12))
            }
        }
        .frame(width: size, height: size)
        .clipShape(RoundedRectangle(cornerRadius: size * 0.24))
        .accessibilityHidden(true)
    }
}

struct CrumbEmptyState: View {
    let title: LocalizedStringKey
    let message: LocalizedStringKey
    let icon: String

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            CrumbArtwork(symbol: icon, size: 44)
            Text(title).font(.headline)
            Text(message).font(.subheadline).foregroundStyle(.secondary)
        }
        .padding(.vertical, 10)
    }
}

struct CrumbActionStyle: ButtonStyle {
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @Environment(\.isEnabled) private var isEnabled
    var primary = true

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.subheadline.weight(.semibold))
            .multilineTextAlignment(.center)
            .padding(.horizontal, 16)
            .padding(.vertical, 9)
            .frame(maxWidth: .infinity, minHeight: 44)
            .foregroundStyle(primary ? Color(uiColor: .systemBackground) : Color.primary)
            .background(primary ? Color.primary : Color.accentColor.opacity(0.12), in: Capsule())
            .opacity(isEnabled ? (configuration.isPressed ? 0.82 : 1) : 0.45)
            .scaleEffect(configuration.isPressed && !reduceMotion ? 0.98 : 1)
            .animation(reduceMotion ? nil : .easeOut(duration: 0.16), value: configuration.isPressed)
    }
}

struct CrumbProgressMark: View {
    var progress: Double?
    var complete = false
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        ZStack {
            Circle().strokeBorder(Color.accentColor.opacity(0.18), lineWidth: 3)
            if complete {
                Image(systemName: "checkmark").font(.headline.weight(.semibold))
                    .transition(reduceMotion ? .identity : .scale.combined(with: .opacity))
            } else if let progress {
                Circle()
                    .trim(from: 0, to: min(max(progress, 0), 1))
                    .stroke(style: StrokeStyle(lineWidth: 3, lineCap: .round))
                    .rotationEffect(.degrees(-90))
                    .padding(1.5)
                    .animation(reduceMotion ? nil : .linear(duration: 0.2), value: progress)
            } else if reduceMotion {
                Image(systemName: "hourglass").font(.subheadline)
            } else {
                ProgressView()
            }
        }
        .foregroundStyle(Color.accentColor)
        .frame(width: 42, height: 42)
        .animation(reduceMotion ? nil : .spring(response: 0.32, dampingFraction: 0.83), value: complete)
        .accessibilityHidden(true)
    }
}

enum CrumbAppearance: String, CaseIterable {
    case system, light, dark

    var colorScheme: ColorScheme? {
        switch self {
        case .system: nil
        case .light: .light
        case .dark: .dark
        }
    }
}

struct CrumbScrollRequest: Equatable {
    var tab: AppTab = .search
    var count = 0
}

struct CrumbScrollRequestKey: EnvironmentKey {
    static let defaultValue = CrumbScrollRequest()
}

extension EnvironmentValues {
    var crumbScrollRequest: CrumbScrollRequest {
        get { self[CrumbScrollRequestKey.self] }
        set { self[CrumbScrollRequestKey.self] = newValue }
    }
}

extension View {
    func crumbPrimaryPage(_ tab: AppTab) -> some View {
        modifier(CrumbPrimaryPage(tab: tab))
    }
}

private struct CrumbPrimaryPage: ViewModifier {
    let tab: AppTab
    @Environment(\.crumbScrollRequest) private var request
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    func body(content: Content) -> some View {
        ScrollViewReader { proxy in
            content
                .onChange(of: request) { value in
                    guard value.tab == tab else { return }
                    withAnimation(reduceMotion ? nil : .easeOut(duration: 0.24)) {
                        proxy.scrollTo(tab, anchor: .top)
                    }
                }
        }
    }
}
