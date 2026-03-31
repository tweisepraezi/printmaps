#!/bin/sh

# ------------------------------------
# Purpose:
# - Build binary for Linux target system.
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
gosec -exclude=G101,G204,G304 ./...

# show compiler version
go version

# -----------------------------------------
# BUILD LINUX
# -----------------------------------------

# linux
env GOOS=linux GOARCH=amd64 go build -v -o build/linux-amd64/printmaps_buildservice

