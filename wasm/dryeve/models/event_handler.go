// ===============================================================
// File: event_handler.go
// Description: Declaration of event handler function types
// Author: DryBearr
// ===============================================================

package models

// SizeChangeHandler handles window resize events.
type SizeChangeHandler func(width int, height int) error

// MouseClickHandler handles mouse click events.
type MouseClickHandler func(point Point2D) error

// MouseDragHandler handles mouse drag move events.
type MouseDragHandler func(point Point2D) error

// MouseDragEndHandler handles the end of a mouse drag.
type MouseDragEndHandler func(point Point2D) error

// KeyDownHandler handles key press events.
type KeyDownHandler func(key Key) error

// SwipeHandler handles swipe direction events.
type SwipeHandler func(direction SwipeDirection) error
