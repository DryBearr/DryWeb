# ===============================================================
# File: build.ps1
# Description: Builds the wasm and serve Windows edition
# Author: DryBearr
# ===============================================================

# Build wasm
cd wasm
$env:GOOS = "js"
$env:GOARCH = "wasm"
go build -o ../static/wasm/game_of_life.wasm ./game_of_life/main.go
go build -o ../static/wasm/snake.wasm ./snake/main.go

# Ensure bin directory exists
New-Item -ItemType Directory -Force -Path ../bin | Out-Null

# Clear environment so we don't mess up the server build
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

# Build server
cd ../serve
go build -o ../bin/serve.exe

cd ..
