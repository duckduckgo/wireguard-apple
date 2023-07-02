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
            url: "https://github.com/DuckDuckGo/wireguard-apple/releases/download/1.1.0/WireGuard.xcframework.zip",
            checksum: "778df46d1523a90a526fb3c13e25c4569287a9e3459896f056f1246beddc4a0d"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
