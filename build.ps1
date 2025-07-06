# ===============================================================
# File: build.ps1
# Description: Builds the wasm and optionally the server binary
# Author: DryBearr
# ===============================================================

param (
    [switch]$NoServe
)

# Build wasm
Set-Location wasm
$env:GOOS = "js"
$env:GOARCH = "wasm"
go build -o ../static/wasm/game_of_life.wasm ./game_of_life/main.go
go build -o ../static/wasm/snake.wasm ./snake/main.go

# Ensure bin directory exists
New-Item -ItemType Directory -Force -Path ../bin | Out-Null

# Clear environment
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

if (-not $NoServe) {
    # Build server
    Set-Location ../serve
    go build -o ../bin/serve.exe
}

Set-Location ..
