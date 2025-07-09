// ===============================================================
// File: render.go
// Description: Provides interface for rendering
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package render

// TODO: godoc
type Size struct {
	Width  int
	Height int
}

// TODO: godoc
type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// TODO: godoc
type SizeChangeHandler func(size Size) error

// TODO: godoc
type Renderer interface {
	RegisterResizeEventListener(handler SizeChangeHandler) error

	DrawFrame(frame *[]Pixel, size Size) error

	GetSize() Size
}
