// ===============================================================
// File: render.go
// Description: Defines godoc for render package and Renderer interface
// Author: DryBearr
// ===============================================================

// TODO: godoc
package render

import "wasm/dryeve/models"

// TODO: godoc
type Renderer interface {
	RenderRect(rect models.Rect, pixel models.Pixel) error
	RenderCircle(circle models.Circle, pixel models.Pixel) error
	RenderPixel(point models.Point2D, pixel models.Pixel) error
	RenderFrame(frame models.RenderFrame) error
	RenderLine(line models.Line, pixel models.Pixel) error
}
