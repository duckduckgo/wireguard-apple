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
            url: "https://github.com/duckduckgo/wireguard-apple/releases/download/1.1.4-relay-poc/WireGuard.xcframework.zip",
            checksum: "015f46fa0389ad1c91f8b2bfc2550dd71cc8eaa87d4c5f94fd41dbd677efad09"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
