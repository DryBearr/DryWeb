// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"wasm/game_of_life/core"
	"wasm/webrender"
)

func main() {
	core.StartGame(webrender.Api)
	select {}
}
