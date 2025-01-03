#!/bin/bash

PROJECT_NAME="ethcracker"

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o ${PROJECT_NAME}-linux
if [ $? -eq 0 ]; then
    echo "Successfully built for Linux."
else
    echo "Failed to build for Linux." >&2
    exit 1
fi

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ${PROJECT_NAME}-macos
if [ $? -eq 0 ]; then
    echo "Successfully built for macOS (Intel)."
else
    echo "Failed to build for macOS (Intel)." >&2
    exit 1
fi

# Build for macOS (ARM, e.g., M1/M2)
GOOS=darwin GOARCH=arm64 go build -o ${PROJECT_NAME}-macos-arm64
if [ $? -eq 0 ]; then
    echo "Successfully built for macOS (ARM)."
else
    echo "Failed to build for macOS (ARM)." >&2
    exit 1
fi

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o ${PROJECT_NAME}-windows.exe
if [ $? -eq 0 ]; then
    echo "Successfully built for Windows."
else
    echo "Failed to build for Windows." >&2
    exit 1
fi

# Summary
echo "Builds completed:"
ls -lh ${PROJECT_NAME}-*