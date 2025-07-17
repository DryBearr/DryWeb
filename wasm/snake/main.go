// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"wasm/snake/core"
	"wasm/webrender"
)

func main() {
	core.StartGame(webrender.Api)
}
