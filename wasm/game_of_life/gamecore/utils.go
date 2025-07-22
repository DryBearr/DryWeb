// ===============================================================
// File: utils.go
// Description: Provides utility functions for the game of life
// Author: DryBearr
// ===============================================================

package gamecore

func Abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}
