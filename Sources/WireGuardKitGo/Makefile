# SPDX-License-Identifier: MIT
#
# Copyright (C) 2018-2019 Jason A. Donenfeld <Jason@zx2c4.com>. All Rights Reserved.

# These are generally passed to us by xcode, but we set working defaults for standalone compilation too.
MACOS_ARCHS ?= x86_64 arm64
MACOS_PLATFORM_NAME ?= macosx
IOS_PLATFORM_NAME ?= iphoneos
IOS_SIMULATOR_PLATFORM_NAME ?= iphonesimulator

IOS_ARCH_PLATFORM_PAIRS := arm64-iphoneos
IOS_SIM_ARCH_PLATFORM_PAIRS := arm64-iphonesimulator x86_64-iphonesimulator

FRAMEWORK_NAME := WireGuard.framework

MACOS_LIBRARY_NAME ?= libwg-go-macos.a
MACOS_FRAMEWORK_DESTINATION ?= macOS/$(FRAMEWORK_NAME)

IOS_BINARY_NAME ?= libwg-go-ios.a
IOS_FRAMEWORK_DESTINATION ?= iOS/$(FRAMEWORK_NAME)

IOS_SIM_BINARY_NAME ?= libwg-go-ios-sim.a
IOS_SIM_FRAMEWORK_DESTINATION ?= Simulator/$(FRAMEWORK_NAME)

MACOS_SDKROOT ?= $(shell xcrun --sdk $(MACOS_PLATFORM_NAME) --show-sdk-path)
IOS_SDKROOT ?= $(shell xcrun --sdk $(IOS_PLATFORM_NAME) --show-sdk-path)
IOS_SIMULATOR_SDKROOT ?= $(shell xcrun --sdk $(IOS_SIMULATOR_PLATFORM_NAME) --show-sdk-path)

CONFIGURATION_BUILD_DIR ?= $(CURDIR)/out
CONFIGURATION_TEMP_DIR ?= $(CURDIR)/.tmp

export PATH := $(PATH):/usr/local/bin:/opt/homebrew/bin
export CC ?= clang

LIPO ?= lipo
DESTDIR ?= $(CONFIGURATION_BUILD_DIR)
BUILDDIR ?= $(CONFIGURATION_TEMP_DIR)/wireguard-go-bridge

MACOS_CFLAGS_PREFIX := -mmacosx-version-min=10.15 -isysroot $(MACOS_SDKROOT) -arch
IOS_CFLAGS_PREFIX := -miphoneos-version-min=14.0 -isysroot $(IOS_SDKROOT) -arch
IOS_SIMULATOR_CFLAGS_PREFIX := -miphonesimulator-version-min=14.0 -isysroot $(IOS_SIMULATOR_SDKROOT) -arch

GOARCH_arm64 := arm64
GOARCH_x86_64 := amd64
GOOS_macosx := darwin
GOOS_iphoneos := ios
GOOS_iphonesimulator := ios

LOG_STATUS_COLOR := $(shell tput setaf 2; tput bold)
LOG_ERROR_COLOR := $(shell tput setaf 1; tput bold)
RESET := $(shell tput sgr0)

build: $(DESTDIR)/$(MACOS_LIBRARY_NAME) $(DESTDIR)/$(IOS_BINARY_NAME) $(DESTDIR)/$(IOS_SIM_BINARY_NAME) $(DESTDIR)/WireGuard.xcframework
build-macos-framework: $(DESTDIR)/macOS/$(FRAMEWORK_NAME)

version-header: $(DESTDIR)/wireguard-go-version.h
xcframework: $(DESTDIR)/WireGuard.xcframework

REAL_GOROOT := $(shell go env GOROOT 2>/dev/null)
export GOROOT := $(BUILDDIR)/goroot
$(GOROOT)/.prepared:
	[ -n "$(REAL_GOROOT)" ]
	mkdir -p "$(GOROOT)"
	rsync -a --delete --exclude=pkg/obj/go-build "$(REAL_GOROOT)/" "$(GOROOT)/"
	cat goruntime-*.diff | patch -p1 -f -N -r- -d "$(GOROOT)"
	touch "$@"

# Compile static library for platform + architecture:

define libwg-go-macos-a
$(BUILDDIR)/libwg-go-macos-$(1).a: export CGO_ENABLED := 1
$(BUILDDIR)/libwg-go-macos-$(1).a: export CGO_CFLAGS := $(MACOS_CFLAGS_PREFIX) $(ARCH)
$(BUILDDIR)/libwg-go-macos-$(1).a: export CGO_LDFLAGS := $(MACOS_CFLAGS_PREFIX) $(ARCH)
$(BUILDDIR)/libwg-go-macos-$(1).a: export GOOS := $(GOOS_$(MACOS_PLATFORM_NAME))
$(BUILDDIR)/libwg-go-macos-$(1).a: export GOARCH := $(GOARCH_$(1))
$(BUILDDIR)/libwg-go-macos-$(1).a: $(GOROOT)/.prepared go.mod
	@echo "$(LOG_STATUS_COLOR)Creating libwg-go-macos-$(1)$(RESET)"
	go build -ldflags=-w -trimpath -v -o "$(BUILDDIR)/libwg-go-macos-$(1).a" -buildmode c-archive
endef
$(foreach ARCH,$(MACOS_ARCHS),$(eval $(call libwg-go-macos-a,$(ARCH))))

define libwg-go-ios-a
$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a: export CGO_ENABLED := 1
$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a: export CGO_CFLAGS := $(IOS_CFLAGS_PREFIX) $(word 1,$(subst -, ,$(1)))
$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a: export CGO_LDFLAGS := $(IOS_CFLAGS_PREFIX) $(word 1,$(subst -, ,$(1)))
$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a: export GOOS := $(GOOS_$(word 2,$(subst -, ,$(1))))
$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a: export GOARCH := $(GOARCH_$(word 1,$(subst -, ,$(1))))
$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a: $(GOROOT)/.prepared go.mod
	@echo "$(LOG_STATUS_COLOR)Creating libwg-go-ios-$(word 1,$(subst -, ,$(1))) with platform $(word 2,$(subst -, ,$(1)))$(RESET)"
	go build -ldflags=-w -trimpath -v -o "$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(1))).a" -buildmode c-archive
endef
$(foreach PAIR,$(IOS_ARCH_PLATFORM_PAIRS),$(eval $(call libwg-go-ios-a,$(PAIR))))

define libwg-go-ios-sim-a
$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a: export CGO_ENABLED := 1
$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a: export CGO_CFLAGS := $(IOS_SIMULATOR_CFLAGS_PREFIX) $(word 1,$(subst -, ,$(1)))
$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a: export CGO_LDFLAGS := $(IOS_SIMULATOR_CFLAGS_PREFIX) $(word 1,$(subst -, ,$(1)))
$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a: export GOOS := $(GOOS_$(word 2,$(subst -, ,$(1))))
$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a: export GOARCH := $(GOARCH_$(word 1,$(subst -, ,$(1))))
$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a: $(GOROOT)/.prepared go.mod
	@echo "$(LOG_STATUS_COLOR)Creating libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))) with platform $(word 2,$(subst -, ,$(1)))$(RESET)"
	go build -ldflags=-w -trimpath -v -o "$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(1))).a" -buildmode c-archive
endef
$(foreach PAIR,$(IOS_SIM_ARCH_PLATFORM_PAIRS),$(eval $(call libwg-go-ios-sim-a,$(PAIR))))

$(DESTDIR)/wireguard-go-version.h: go.mod $(GOROOT)/.prepared
	 sed -E -n 's/.*golang\.zx2c4\.com\/wireguard +v[0-9.]+-[0-9]+-([0-9a-f]{8})[0-9a-f]{4}.*/#define WIREGUARD_GO_VERSION "\1"/p' "$<" > "$@"

# Combine macOS slices into binary with lipo:

$(DESTDIR)/$(MACOS_LIBRARY_NAME): $(foreach ARCH,$(MACOS_ARCHS),$(BUILDDIR)/libwg-go-macos-$(ARCH).a)
	@echo "$(LOG_STATUS_COLOR)Combining macOS architectures from $^ into $@$(RESET)"
	@mkdir -vp "$(DESTDIR)"
	$(LIPO) -create -output "$@" $^

$(DESTDIR)/$(IOS_BINARY_NAME): $(foreach PAIR,$(IOS_ARCH_PLATFORM_PAIRS),$(BUILDDIR)/libwg-go-ios-$(word 1,$(subst -, ,$(PAIR))).a)
	@echo "$(LOG_STATUS_COLOR)Combining iOS device architectures from $^ into $@$(RESET)"
	@mkdir -vp "$(DESTDIR)"
	@cp $^ $@

$(DESTDIR)/$(IOS_SIM_BINARY_NAME): $(foreach PAIR,$(IOS_SIM_ARCH_PLATFORM_PAIRS),$(BUILDDIR)/libwg-go-ios-sim-$(word 1,$(subst -, ,$(PAIR))).a)
	@echo "$(LOG_STATUS_COLOR)Combining simulator architectures from $^ into $@$(RESET)"
	@mkdir -vp "$(DESTDIR)"
	$(LIPO) -create -output "$@" $^

# Create frameworks:

$(DESTDIR)/macOS/$(FRAMEWORK_NAME): $(DESTDIR)/$(MACOS_LIBRARY_NAME)
	@echo "$(LOG_STATUS_COLOR)Building macOS framework: $@$(RESET)"

	@mkdir -vp $@/Versions
	@mkdir -vp $@/Versions/A
	@mkdir -vp $@/Versions/A/Headers
	@mkdir -vp $@/Versions/A/Modules
	@mkdir -vp $@/Versions/A/Resources

	@cp Info.plist $@/Versions/A/Resources
	@cp include-framework/wireguard.h $@/Versions/A/Headers
	@cp include-framework/module.modulemap $@/Versions/A/Modules/module.modulemap
	@cp $^ $@/Versions/A/WireGuard

	@cd out/macOS/$(FRAMEWORK_NAME)/Versions && ln -sfhv A Current
	@cd out/macOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Headers Headers
	@cd out/macOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Modules Modules
	@cd out/macOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Resources Resources
	@cd out/macOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/WireGuard WireGuard

$(DESTDIR)/iOS/$(FRAMEWORK_NAME): $(DESTDIR)/$(IOS_BINARY_NAME)
	@echo "$(LOG_STATUS_COLOR)Building iOS framework: $@$(RESET)"

	@mkdir -vp $@/Versions
	@mkdir -vp $@/Versions/A
	@mkdir -vp $@/Versions/A/Headers
	@mkdir -vp $@/Versions/A/Modules
	@mkdir -vp $@/Versions/A/Resources

	@cp Info.plist $@/Versions/A/Resources
	@cp include-framework/wireguard.h $@/Versions/A/Headers
	@cp include-framework/module.modulemap $@/Versions/A/Modules/module.modulemap
	@cp $^ $@/Versions/A/WireGuard

	@cd out/iOS/$(FRAMEWORK_NAME)/Versions && ln -sfhv A Current
	@cd out/iOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Headers Headers
	@cd out/iOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Modules Modules
	@cd out/iOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Resources Resources
	@cd out/iOS/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/WireGuard WireGuard

$(DESTDIR)/Simulator/$(FRAMEWORK_NAME): $(DESTDIR)/$(IOS_SIM_BINARY_NAME)
	@echo "$(LOG_STATUS_COLOR)Building iOS simulator framework: $@$(RESET)"

	@mkdir -vp $@/Versions
	@mkdir -vp $@/Versions/A
	@mkdir -vp $@/Versions/A/Headers
	@mkdir -vp $@/Versions/A/Modules
	@mkdir -vp $@/Versions/A/Resources

	@cp Info.plist $@/Versions/A/Resources
	@cp include-framework/wireguard.h $@/Versions/A/Headers
	@cp include-framework/module.modulemap $@/Versions/A/Modules/module.modulemap
	@cp $^ $@/Versions/A/WireGuard

	@cd out/Simulator/$(FRAMEWORK_NAME)/Versions && ln -sfhv A Current
	@cd out/Simulator/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Headers Headers
	@cd out/Simulator/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Modules Modules
	@cd out/Simulator/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/Resources Resources
	@cd out/Simulator/$(FRAMEWORK_NAME) && ln -sfhv Versions/Current/WireGuard WireGuard

# Zip it all up in an XCFramework:

$(DESTDIR)/WireGuard.xcframework: $(DESTDIR)/macOS/$(FRAMEWORK_NAME) $(DESTDIR)/iOS/$(FRAMEWORK_NAME) $(DESTDIR)/Simulator/$(FRAMEWORK_NAME)
	@echo "$(LOG_STATUS_COLOR)Creating XCFramework$(RESET)"
	xcodebuild -create-xcframework -framework $(DESTDIR)/macOS/$(FRAMEWORK_NAME) -framework $(DESTDIR)/iOS/$(FRAMEWORK_NAME) -framework $(DESTDIR)/Simulator/$(FRAMEWORK_NAME) -output $(DESTDIR)/WireGuard.xcframework

# Useful actions:

clean:
	@echo "$(LOG_STATUS_COLOR)Cleaning build artifacts$(RESET)"
	rm -rf "$(BUILDDIR)" "$(DESTDIR)"

install: build

.PHONY: build clean install version-header xcframework
