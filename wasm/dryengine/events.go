// ===============================================================
// File: events.go
// Description: Provides interface for event handling
// Author: DryBearr
// ===============================================================

package dryengine

// Key represents a keyboard key.
type Key string

const (
	WKey Key = "w"
	AKey Key = "a"
	SKey Key = "s"
	DKey Key = "d"
	PKey Key = "p"
	RKey Key = "r"
)

// Coordinate2D represents a 2D point with X and Y values.
type Coordinate2D struct {
	X int
	Y int
}

// SwipeDirection represents a direction of swipe in 2D space.
type SwipeDirection Coordinate2D

var (
	SwipeLeft  = SwipeDirection{X: -1, Y: 0}
	SwipeRight = SwipeDirection{X: 1, Y: 0}
	SwipeUp    = SwipeDirection{X: 0, Y: -1}
	SwipeDown  = SwipeDirection{X: 0, Y: 1}
)

// SizeChangeHandler handles window resize events.
type SizeChangeHandler func(s Size) error

// MouseClickHandler handles mouse click events.
type MouseClickHandler func(c Coordinate2D) error

// MouseDragHandler handles mouse drag events.
type MouseDragHandler func(c Coordinate2D) error

// MouseDragEndHandler handles the end of a mouse drag event.
type MouseDragEndHandler func(c Coordinate2D) error

// KeyDownHandler handles key down events.
type KeyDownHandler func(key Key) error

// TODO: KeyUpHandler handles key up events.
// type KeyUpHandler func(key Key) error

// SwipeHandler handles swipe direction events.
type SwipeHandler func(direction SwipeDirection) error

// DryEvents defines methods for registering various input and window event handlers.
type DryEvents interface {
	RegisterResizeEventListener(handler SizeChangeHandler) error
	RegisterMouseClickEventListener(handler MouseClickHandler) error
	RegisterMouseDragEventListener(handler MouseDragHandler) error
	RegisterMouseDragEndEventListener(handler MouseDragEndHandler) error
	RegisterKeyDownEventListener(handler KeyDownHandler) error
	RegisterSwipeEventListener(handler SwipeHandler) error
}
