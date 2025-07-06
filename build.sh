#!/usr/bin/env bash

# ===============================================================
# File: build.sh
# Description: Builds the wasm and optionally the server binary
# Author: DryBearr
# ===============================================================

build_serve=true

for arg in "$@"; do
  if [ "$arg" = "--no-serve" ]; then
    build_serve=false
  fi
done

cd wasm
GOOS=js GOARCH=wasm go build -o ../static/wasm/game_of_life.wasm ./game_of_life/main.go
GOOS=js GOARCH=wasm go build -o ../static/wasm/snake.wasm ./snake/main.go
cd ..

if [ "$build_serve" = true ]; then
  cd serve
  mkdir -p ../bin
  go build -o ../bin/serve
  cd ..
fi
