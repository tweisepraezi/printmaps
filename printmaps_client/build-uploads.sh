#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds uploads (tar.gz or zip) for printmaps commandline client.
# 
# Releases:
# - v1.0.0 - 2026-03-31: initial release
# ------------------------------------

# set -o xtrace
set -o verbose

# recreate directory
rm -rf ./uploads
mkdir -p ./uploads

# uploads 'linux'
tar -cvzf ./uploads/printmaps_linux_amd64.tar.gz -C ./build/linux-amd64 printmaps
tar -cvzf ./uploads/printmaps_linux_arm64.tar.gz -C ./build/linux-arm64 printmaps

# uploads 'darwin' (macOS)
tar -cvzf ./uploads/printmaps_darwin_amd64.tar.gz -C ./build/darwin-amd64 printmaps
tar -cvzf ./uploads/printmaps_darwin_arm64.tar.gz -C ./build/darwin-arm64 printmaps

# uploads 'windows'
zip -j ./uploads/printmaps_windows_amd64.zip ./build/windows-amd64/printmaps.exe
zip -j ./uploads/printmaps_windows_arm64.zip ./build/windows-arm64/printmaps.exe
