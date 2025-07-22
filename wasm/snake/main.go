// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"time"
	"wasm/dryeve/engine"
	"wasm/dryeve/web"
	"wasm/snake/gamecore"
)

func main() {
	gameEvents := web.NewWebEvents()
	gameRenderer := web.NewWebRenderer()
	gameEngine := engine.NewEngine(gameRenderer, gameEvents, 16*time.Millisecond, 1000)

	gamecore.StartGame(*gameEngine)
}
