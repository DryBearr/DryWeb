// ===============================================================
// File: render.go
// Description: Provides interface for rendering
// Author: DryBearr
// ===============================================================

package render

type Key string

const (
	WKey Key = "w"
	AKey Key = "a"
	SKey Key = "s"
	DKey Key = "d"
	PKey Key = "p"
	RKey Key = "r"
)

// TODO: godoc
type Coordinate struct {
	X int
	Y int
}

type SwipeDirection Coordinate

var (
	SwipeLeft = SwipeDirection{
		X: -1,
		Y: 0,
	}

	SwipeRight = SwipeDirection{
		X: 1,
		Y: 0,
	}

	SwipeUp = SwipeDirection{
		X: 0,
		Y: -1,
	}

	SwipeDown = SwipeDirection{
		X: 0,
		Y: 1,
	}
)

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
type MouseDragEndHandler func(c Coordinate) error

// Key event handlers
type KeyDownHandler func(key Key) error

//TODO: type KeyUpHandler func(key Key) error

// Swipe events
type SwipeHandler func(direction SwipeDirection) error

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
	RegisterMouseDragEndEventListener(handler MouseDragEndHandler) error

	//Key Events
	RegisterKeyDownEventListener(handler KeyDownHandler) error

	//Swipe Events
	RegisterSwipeEventListener(handler SwipeHandler) error
}
