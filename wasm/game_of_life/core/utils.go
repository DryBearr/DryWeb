// ===============================================================
// File: drawing.go
// Description: Provides utility functions for the game of life
// Author: DryBearr
// ===============================================================

package core

func abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}
