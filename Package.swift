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
            url: "https://github.com/duckduckgo/wireguard-apple/releases/download/1.1.4-relay-poc-v4/WireGuard.xcframework.zip",
            checksum: "37cf93ac8cb05b2aacac93f8ab38f262f8a27022fc6e972a17238c1b6e44e1e4"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
