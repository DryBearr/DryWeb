// ===============================================================
// File: key.go
// Description: Defines keyboard key model
// Author: DryBearr
// ===============================================================

package models

//TODO: add all keys and godoc

type Key int

const (
	KeyUnknown Key = iota

	KeyA
	KeyW
	KeyS
	KeyD

	KeyLeft
	KeyRight
	KeyUp
	KeyDown
)

func (k Key) String() string {
	switch k {
	case KeyA:
		return "A"
	case KeyW:
		return "W"
	case KeyS:
		return "S"
	case KeyD:
		return "D"
	case KeyLeft:
		return "Left"
	case KeyRight:
		return "Right"
	case KeyUp:
		return "Up"
	case KeyDown:
		return "Down"
	default:
		return "Unknown"
	}
}
