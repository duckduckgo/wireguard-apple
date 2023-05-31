# SPDX-License-Identifier: MIT
#
# Copyright (C) 2018-2019 Jason A. Donenfeld <Jason@zx2c4.com>. All Rights Reserved.

# These are generally passed to us by xcode, but we set working defaults for standalone compilation too.
MACOS_ARCHS ?= x86_64 arm64
MACOS_PLATFORM_NAME ?= macosx

IOS_ARCHS ?= arm64
IOS_PLATFORM_NAME ?= iphoneos
IPHONESIMULATOR_PLATFORM_NAME ?= iphonesimulator

MACOS_LIBRARY_NAME ?= libwg-go-macos.a
MACOS_FRAMEWORK_DESTINATION ?= macOS/WireGuard.framework
MACOS_FRAMEWORK_NAME ?= WireGuard.framework

IOS_BINARY_NAME ?= libwg-go-ios.a

MACOS_SDKROOT ?= $(shell xcrun --sdk $(MACOS_PLATFORM_NAME) --show-sdk-path)
IOS_SDKROOT ?= $(shell xcrun --sdk $(IOS_PLATFORM_NAME) --show-sdk-path)

CONFIGURATION_BUILD_DIR ?= $(CURDIR)/out
CONFIGURATION_TEMP_DIR ?= $(CURDIR)/.tmp

export PATH := $(PATH):/usr/local/bin:/opt/homebrew/bin
export CC ?= clang
LIPO ?= lipo
DESTDIR ?= $(CONFIGURATION_BUILD_DIR)
BUILDDIR ?= $(CONFIGURATION_TEMP_DIR)/wireguard-go-bridge

CFLAGS_PREFIX := -mmacosx-version-min=10.15 -isysroot $(MACOS_SDKROOT) -arch
IOS_CFLAGS_PREFIX := $(if $(DEPLOYMENT_TARGET_CLANG_FLAG_NAME),-$(DEPLOYMENT_TARGET_CLANG_FLAG_NAME)=$($(DEPLOYMENT_TARGET_CLANG_ENV_NAME)),) -isysroot $(IOS_SDKROOT) -arch
GOARCH_arm64 := arm64
GOARCH_x86_64 := amd64
GOOS_macosx := darwin
GOOS_iphoneos := ios

LOG_STATUS_COLOR := $(shell tput setaf 2; tput bold)
LOG_ERROR_COLOR := $(shell tput setaf 1; tput bold)
RESET := $(shell tput sgr0)

build: $(DESTDIR)/$(MACOS_LIBRARY_NAME) $(DESTDIR)/$(IOS_BINARY_NAME) $(DESTDIR)/WireGuard.xcframework
build-framework: $(DESTDIR)/macOS/$(MACOS_FRAMEWORK_NAME)

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
$(BUILDDIR)/libwg-go-macos-$(1).a: export CGO_CFLAGS := $(CFLAGS_PREFIX) $(ARCH)
$(BUILDDIR)/libwg-go-macos-$(1).a: export CGO_LDFLAGS := $(CFLAGS_PREFIX) $(ARCH)
$(BUILDDIR)/libwg-go-macos-$(1).a: export GOOS := $(GOOS_$(MACOS_PLATFORM_NAME))
$(BUILDDIR)/libwg-go-macos-$(1).a: export GOARCH := $(GOARCH_$(1))
$(BUILDDIR)/libwg-go-macos-$(1).a: $(GOROOT)/.prepared go.mod
	@echo "$(LOG_STATUS_COLOR)Creating libwg-go-macos-$(1)$(RESET)"
	go build -ldflags=-w -trimpath -v -o "$(BUILDDIR)/libwg-go-macos-$(1).a" -buildmode c-archive
endef
$(foreach ARCH,$(MACOS_ARCHS),$(eval $(call libwg-go-macos-a,$(ARCH))))

define libwg-go-ios-a
$(BUILDDIR)/libwg-go-ios-$(1).a: export CGO_ENABLED := 1
$(BUILDDIR)/libwg-go-ios-$(1).a: export CGO_CFLAGS := $(IOS_CFLAGS_PREFIX) $(ARCH)
$(BUILDDIR)/libwg-go-ios-$(1).a: export CGO_LDFLAGS := $(IOS_CFLAGS_PREFIX) $(ARCH)
$(BUILDDIR)/libwg-go-ios-$(1).a: export GOOS := $(GOOS_$(IOS_PLATFORM_NAME))
$(BUILDDIR)/libwg-go-ios-$(1).a: export GOARCH := $(GOARCH_$(1))
$(BUILDDIR)/libwg-go-ios-$(1).a: $(GOROOT)/.prepared go.mod
	@echo "$(LOG_STATUS_COLOR)Creating libwg-go-ios-$(1)$(RESET)"
	go build -ldflags=-w -trimpath -v -o "$(BUILDDIR)/libwg-go-ios-$(1).a" -buildmode c-archive
endef
$(foreach ARCH,$(IOS_ARCHS),$(eval $(call libwg-go-ios-a,$(ARCH))))

$(DESTDIR)/wireguard-go-version.h: go.mod $(GOROOT)/.prepared
	 sed -E -n 's/.*golang\.zx2c4\.com\/wireguard +v[0-9.]+-[0-9]+-([0-9a-f]{8})[0-9a-f]{4}.*/#define WIREGUARD_GO_VERSION "\1"/p' "$<" > "$@"

# Combine macOS slices into binary with lipo:

$(DESTDIR)/$(MACOS_LIBRARY_NAME): $(foreach ARCH,$(MACOS_ARCHS),$(BUILDDIR)/libwg-go-macos-$(ARCH).a)
	@echo "$(LOG_STATUS_COLOR)Combining macOS architectures from $^ into $@$(RESET)"
	@mkdir -vp "$(DESTDIR)"
	$(LIPO) -create -output "$@" $^

$(DESTDIR)/$(IOS_BINARY_NAME): $(foreach ARCH,$(IOS_ARCHS),$(BUILDDIR)/libwg-go-ios-$(ARCH).a)
	@echo "$(LOG_STATUS_COLOR)Combining architectures into $@$(RESET)"
	@mkdir -vp "$(DESTDIR)"
	@cp $^ $@

$(DESTDIR)/macOS/$(MACOS_FRAMEWORK_NAME): $(DESTDIR)/$(MACOS_LIBRARY_NAME)
	@echo "$(LOG_STATUS_COLOR)Building framework: $@$(RESET)"

	@mkdir -vp "$(DESTDIR)/macOS/$(MACOS_FRAMEWORK_NAME)/Headers"
	@mkdir -vp "$(DESTDIR)/macOS/$(MACOS_FRAMEWORK_NAME)/Modules"

	@cp Info.plist $@
	@cp include/wireguard.h $@/Headers
	@cp include/module.modulemap $@/Modules
	@cp $^ $@/WireGuard

$(DESTDIR)/WireGuard.xcframework: $(DESTDIR)/macOS/$(MACOS_FRAMEWORK_NAME)
	@echo "$(LOG_STATUS_COLOR)Creating XCFramework$(RESET)"
	xcodebuild -create-xcframework -framework $(DESTDIR)/macOS/$(MACOS_FRAMEWORK_NAME) -output $(DESTDIR)/WireGuard.xcframework

clean:
	@echo "$(LOG_STATUS_COLOR)Cleaning build artifacts$(RESET)"
	rm -rf "$(BUILDDIR)" "$(DESTDIR)"

install: build

.PHONY: build build-framework clean install version-header xcframework
