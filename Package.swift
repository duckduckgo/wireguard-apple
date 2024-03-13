// swift-tools-version:5.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "WireGuardKit",
    platforms: [
        .iOS(.v14),
        .macOS(.v10_15),
    ],
    products: [
        .library(name: "WireGuard", targets: ["WireGuard", "_WireGuardDummy"]),
    ],
    targets: [
        .binaryTarget(
            name: "WireGuard",
            url: "https://github.com/DuckDuckGo/wireguard-apple/releases/download/1.1.2/WireGuard.xcframework.zip",
            checksum: "68434186e998f6cd75c0fd840b0b82847eb18737f8bf94eaf9f6bb4b35239dcc"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
