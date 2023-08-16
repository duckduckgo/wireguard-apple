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
            url: "https://github.com/DuckDuckGo/wireguard-apple/releases/download/1.1.1/WireGuard.xcframework.zip",
            checksum: "068099fd1ad6c1ff8c3493d4bf8e5d7610b1f82722c967cf00e5dd6558a25417"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
