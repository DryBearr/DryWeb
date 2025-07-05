#!/usr/bin/env bash

# ===============================================================
# File: build.sh
# Description: Builds the wasm and serve linux edition
# Author: DryBearr
# ===============================================================


cd wasm

# Build wasm
GOOS=js GOARCH=wasm go build -o ../static/wasm/main.wasm .

cd ../serve

# Create bin directory if it doesn't exist
mkdir -p ../bin

# Build server binary into bin/
go build -o ../bin/serve

cd ..
