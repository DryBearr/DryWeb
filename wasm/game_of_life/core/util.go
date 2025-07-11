// ===============================================================
// File: util.go
// Description: Utility functions for the game of life logic
// Author: DryBearr
// ===============================================================

package core

import (
	"sort"
	"wasm/render"
)

func abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}

func BresenhamLine(c0, c1 render.Coordinate) *[]render.Coordinate {
	x0, y0 := c0.X, c0.Y
	x1, y1 := c1.X, c1.Y

	diffX := abs(x0 - x1)
	diffY := abs(y0 - y1)

	stepX := 1
	if x0 > x1 {
		stepX = -1
	}

	stepY := 1
	if y0 > y1 {
		stepY = -1
	}

	err := diffX - diffY

	var coordinates []render.Coordinate

	for {
		coordinates = append(coordinates, render.Coordinate{X: x0, Y: y0})
		if x0 == x1 && y0 == y1 {
			break
		}

		err2 := 2 * err

		if err2 > -diffY {
			err -= diffY
			x0 += stepX
		}

		if err2 < diffX {
			err += diffX
			y0 += stepY
		}
	}

	sort.Slice(coordinates, func(i, j int) bool {
		if coordinates[i].Y == coordinates[j].Y {
			return coordinates[i].X < coordinates[j].X
		}
		return coordinates[i].Y < coordinates[j].Y
	})

	return &coordinates
}
