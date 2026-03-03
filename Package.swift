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
            url: "https://github.com/duckduckgo/wireguard-apple/releases/download/1.1.4-relay-poc-v2/WireGuard.xcframework.zip",
            checksum: "df71ace095083a1259cc7739e2fa0510eb4035cbadbd6d3821d051315c459a69"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
