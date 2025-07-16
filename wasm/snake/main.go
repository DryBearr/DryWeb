// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"wasm/render"
	"wasm/snake/core"
	"wasm/webrender"
)

func main() {
	api := render.Renderer(webrender.Api)

	core.StartGame(api)
}
