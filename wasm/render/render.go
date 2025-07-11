// ===============================================================
// File: render.go
// Description: Provides interface for rendering
// Author: DryBearr
// ===============================================================

package render

// TODO: godoc
type Coordinate struct {
	X int
	Y int
}

// TODO: godoc
type Size struct {
	Width  int
	Height int
}

func (this *Size) EqualOrGreater(other Size) bool {
	return this.Width >= other.Width && this.Height >= other.Height
}

// TODO: godoc
type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Window Event handler
type SizeChangeHandler func(s Size) error

// Mouse event handlers
type MouseClickHandler func(c Coordinate) error
type MouseDragHandler func(c Coordinate) error

// TODO: godoc
type Renderer interface {
	DrawFrame(frame *[][]Pixel, s Size) error
	DrawFramePartly(frame *[][]Pixel, s Size, c Coordinate) error

	GetSize() Size

	//Window Events
	RegisterResizeEventListener(handler SizeChangeHandler) error

	//Mouse Events
	RegisterMouseClickEventListener(handler MouseClickHandler) error
	RegisterMouseDragEventListener(handler MouseDragHandler) error
}
