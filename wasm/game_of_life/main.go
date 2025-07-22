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
	"wasm/game_of_life/gamecore"
)

func main() {
	events := web.NewWebEvents()
	renderer := web.NewWebRenderer()
	gameEngine := engine.NewEngine(renderer, events, 16*time.Millisecond, 1000)

	gamecore.StartGame(*gameEngine)
}
