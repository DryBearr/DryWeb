// ===============================================================
// File: render.go
// Description: Provides interface for rendering
// Author: DryBearr
// ===============================================================

package dryengine

import "time"

// RenderFrame represents a single frame to be rendered.
// It includes the optional top-left coordinate (C) for partial rendering,
// the frame size, and the actual pixel data.
type RenderFrame struct {
	C         *Coordinate2D // Optional top-left coordinate for partial frame
	FrameSize Size          // Size of the frame
	Frame     *[][]Pixel    // 2D array of pixels representing the frame
}

// Size represents the dimensions (width and height) of a renderable area.
type Size struct {
	Width  int
	Height int
}

// EqualOrGreater returns true if the current size is equal to or larger
// than the given size in both width and height.
func (this *Size) EqualOrGreater(other Size) bool {
	return this.Width >= other.Width && this.Height >= other.Height
}

// Pixel represents a single RGBA pixel.
type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// DryRenderer defines the interface for rendering functionality.
// It provides methods to start a render loop, submit frames,
// and query the current render size.
type DryRenderer interface {
	// StartRenderLoop begins the render loop with the specified latency.
	StartRenderLoop(latency time.Duration)

	// AddFrame queues a new frame to be rendered.
	AddFrame(renderFrame *RenderFrame)

	// GetSize returns the current render surface size.
	GetSize() Size
}
