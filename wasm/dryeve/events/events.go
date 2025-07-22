// ===============================================================
// File: events.go
// Description: Defines events interface and package event godoc
// Author: DryBearr
// ===============================================================

// TODO: godoc
package events

import "wasm/dryeve/models"

// Events defines methods for registering various input and window event handlers.
type Events interface {
	RegisterResizeEventListener(handler models.SizeChangeHandler) error
	RegisterMouseClickEventListener(handler models.MouseClickHandler) error
	RegisterMouseDragEventListener(handler models.MouseDragHandler) error
	RegisterMouseDragEndEventListener(handler models.MouseDragEndHandler) error
	RegisterKeyDownEventListener(handler models.KeyDownHandler) error
	RegisterSwipeEventListener(handler models.SwipeHandler) error
}
