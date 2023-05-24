// swift-tools-version:5.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "WireGuardKit",
    platforms: [
        .macOS(.v10_15),
    ],
    products: [
        .library(name: "WireGuard", targets: ["WireGuard", "_WireGuardDummy"]),
    ],
    targets: [
        .binaryTarget(
            name: "WireGuard",
            url: "https://github.com/DuckDuckGo/wireguard-apple/releases/download/1.0.0/WireGuard.xcframework.zip",
            checksum: "06845e28d1d71f53fadb074d999d5141768dcaeba0d0dc3dcc7e2399641bcd3f"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
