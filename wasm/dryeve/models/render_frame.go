// ===============================================================
// File: frame.go
// Description: Defines model for image 2d frame to render
// Author: DryBearr
// ===============================================================

package models

// TODO: godoc
type RenderFrame struct {
	C     *Point2D   // Optional top-left coordinate for partial frame
	Frame *[][]Pixel // 2D array of pixels representing the frame
}
