// ===============================================================
// File: swipe.go
// Description: Defines model for swipe
// Author: DryBearr
// ===============================================================

package models

type SwipeDirection Point2D

var (
	SwipeLeft  = SwipeDirection{X: -1, Y: 0}
	SwipeRight = SwipeDirection{X: 1, Y: 0}
	SwipeUp    = SwipeDirection{X: 0, Y: -1}
	SwipeDown  = SwipeDirection{X: 0, Y: 1}
)
