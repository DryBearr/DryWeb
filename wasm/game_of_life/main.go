// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"syscall/js"
)

const (
	WINDOW_WIDTH  = 500
	WINDOW_HEIGHT = 500

	WINDOW_WIDTH_MB  = 300
	WINDOW_HEIGHT_MB = 300
)

var (
	ctx js.Value
)

func main() {
	//TODO:
	window := js.Global().Get("window")
	viewportWidth := window.Get("innerWidth").Int()
	viewportHeight := window.Get("innerHeight").Int()

	document := js.Global().Get("document")

	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", "window")

	width := WINDOW_WIDTH
	height := WINDOW_HEIGHT

	if viewportHeight < WINDOW_HEIGHT || viewportWidth < WINDOW_WIDTH {
		width = WINDOW_WIDTH_MB
		height = WINDOW_HEIGHT_MB
	}

	canvas.Set("width", width)
	canvas.Set("height", height)

	document.Call("querySelector", "main").Call("append", canvas)
	ctx = canvas.Call("getContext", "2d")
	ctx.Set("fillStyle", "white")
	ctx.Set("strokeStyle", "white")
	ctx.Set("font", "48px serif")
	ctx.Call("fillText", "game of life", 10, 50)

	select {}
}
