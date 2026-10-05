// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "MapleSAP",
    platforms: [.iOS(.v16), .macOS(.v13)],
    products: [.library(name: "MapleSAP", targets: ["MapleSAP"])],
    targets: [
        .target(name: "MemoryProbe", publicHeadersPath: "include"),
        .target(name: "MapleSAP", dependencies: ["MemoryProbe"]),
        .testTarget(name: "MapleSAPTests", dependencies: ["MapleSAP"])
    ]
)
