// ===============================================================
// File: circle.go
// Description: Defines circle
// Author: DryBearr
// ===============================================================

package models

type Circle struct {
	Center Point2D

	R float32

	StartAngle float32
	EndAngle   float32
}
