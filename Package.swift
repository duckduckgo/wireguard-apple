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
            url: "https://github.com/DuckDuckGo/wireguard-apple/releases/download/1.1.3/WireGuard.xcframework.zip",
            checksum: "cd8998e9d9db01484ad19b2a2c425966b4d8116db97a65928fb2d7c87a277ce7"
        ),
        .target(name: "_WireGuardDummy")
    ]

)
