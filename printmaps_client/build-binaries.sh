#!/bin/sh

# ------------------------------------
# Purpose:
# - Build binaries for target systems.
#
# Releases:
# - v1.0.0 - 2026-03-30: initial release
# ------------------------------------

set -o errexit
# set -o verbose

# -----------------------------------------
# CLEANUP & LINTING
# -----------------------------------------

# clean up
rm -rf build
mkdir -p build

# go mod
go mod tidy

# populate vendor directory
go mod vendor

# linting
golangci-lint run --no-config --enable gocritic
revive

# security
govulncheck ./...
gosec -exclude=G304 ./...

# show compiler version
go version

# -----------------------------------------
# BUILD DARWIN (macOS)
# -----------------------------------------

env GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/printmaps
env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/printmaps

# -----------------------------------------
# BUILD LINUX
# -----------------------------------------

env GOOS=linux GOARCH=amd64 go build -v -o build/linux-amd64/printmaps
env GOOS=linux GOARCH=arm64 go build -v -o build/linux-arm64/printmaps

# -----------------------------------------
# BUILD WINDOWS
# -----------------------------------------

env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/printmaps.exe
env GOOS=windows GOARCH=arm64 go build -o build/windows-arm64/printmaps.exe
