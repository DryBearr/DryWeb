#!/usr/bin/env bash

# ===============================================================
# File: build_wasm.sh
# Description: Builds the wasm
# Author: DryBearr
# ===============================================================

set -e

# Resolve script directory (even when run from outside)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
dest_dir="./bin"

show_help() {
  echo "Usage: build_wasm.sh [--output <dir>]"
  echo ""
  echo "Options:"
  echo "  --output, -o DIR  Set output directory (default: ./bin)"
  echo "  --help, -h        Show this help message"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --output|-o)
      dest_dir="$2"
      shift 2
      ;;
    --help|-h)
      show_help
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      show_help
      exit 1
      ;;
  esac
done

mkdir -p "$dest_dir"
dest_dir="$(cd "$dest_dir"; pwd)"
echo "Output directory: $dest_dir"

# Build WASM modules
pushd "$SCRIPT_DIR/wasm" > /dev/null

if [ -f "./game_of_life/main.go" ]; then
  echo "Building game_of_life.wasm..."
  GOOS=js GOARCH=wasm go build -v -o "$dest_dir/game_of_life.wasm" ./game_of_life/main.go || {
    echo "Failed to build game_of_life"
    exit 1
  }
else
  echo "game_of_life/main.go not found. Skipping."
fi

if [ -f "./snake/main.go" ]; then
  echo "Building snake.wasm..."
  GOOS=js GOARCH=wasm go build -v -o "$dest_dir/snake.wasm" ./snake/main.go || {
    echo "Failed to build snake"
    exit 1
  }
else
  echo "snake/main.go not found. Skipping."
fi

popd > /dev/null
