#!/usr/bin/env bash

# ===============================================================
# File: build.sh
# Description: Builds the wasm and serve linux edition
# Author: DryBearr
# ===============================================================


cd wasm

# Build wasm
GOOS=js GOARCH=wasm go build -o ../static/wasm/game_of_life.wasm ./game_of_life/main.go
GOOS=js GOARCH=wasm go build -o ../static/wasm/snake.wasm ./snake/main.go

cd ../serve

# Create bin directory if it doesn't exist
mkdir -p ../bin

# Build server binary into bin/
go build -o ../bin/serve

cd ..
