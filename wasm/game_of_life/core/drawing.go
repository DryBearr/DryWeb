// ===============================================================
// File: drawing.go
// Description: Logic of drawing for the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"time"
	"wasm/render"
)

var (
	coordinateChan = make(chan render.Coordinate, 100)

	color uint8 = 255
	pixel       = render.Pixel{
		R: color,
		G: color,
		B: color,
		A: 255,
	}

	pointFrame = [][]render.Pixel{
		{
			pixel,
		},
	}
	pointSize = render.Size{Width: 1, Height: 1}
)

func StartDrawingLoop() {
	go func() {
		var prev *render.Coordinate
		timeout := 20 * time.Millisecond
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		for {
			select {
			case c, ok := <-coordinateChan:
				if !ok {
					return
				}

				if prev == nil {
					prev = &c
				} else {
					drawLine(*prev, c)
					prev = &c
				}

				// Reset timer
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(timeout)

			case <-timer.C:
				if prev != nil {
					drawPoint(*prev)
				}

				prev = nil

				timer.Reset(timeout)
			}
		}
	}()
}

func AddQueue(c render.Coordinate) {
	coordinateChan <- c
}

func drawLine(start render.Coordinate, end render.Coordinate) {
	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y

	diffX := abs(x0 - x1)
	diffY := abs(y0 - y1)

	minX := min(x0, x1)
	minY := min(y0, y1)

	tempSize := render.Size{
		Width:  diffX + 1,
		Height: diffY + 1,
	}

	tempFrame := make([][]render.Pixel, tempSize.Height)

	frameMutex.Lock()
	for row := range tempFrame {
		tempFrame[row] = make([]render.Pixel, tempSize.Width)

		for column := range tempFrame[row] {
			frameY := minY + row
			frameX := minX + column

			if frameY >= 0 && frameY < len(frame2D) && frameX >= 0 && frameX < len(frame2D[0]) {
				tempFrame[row][column] = frame2D[frameY][frameX]
			}
		}
	}
	frameMutex.Unlock()

	stepX := 1
	if x0 > x1 {
		stepX = -1
	}

	stepY := 1
	if y0 > y1 {
		stepY = -1
	}

	err := diffX - diffY
	frameMutex.Lock()
	for {
		// Compute tempFrame indices
		tempX := x0 - minX
		tempY := y0 - minY

		// Check tempFrame bounds
		if tempY >= 0 && tempY < tempSize.Height && tempX >= 0 && tempX < tempSize.Width {
			tempFrame[tempY][tempX] = pixel
		}

		// Update main frame if needed
		if y0 >= 0 && y0 < len(frame2D) && x0 >= 0 && x0 < len(frame2D[0]) {
			frame2D[y0][x0] = pixel
		}

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
	frameMutex.Unlock()

	api.DrawFramePartly(&tempFrame, tempSize, render.Coordinate{X: minX, Y: minY})
}

func drawPoint(c render.Coordinate) {
	setPixel(pixel, c)

	api.DrawFramePartly(&pointFrame, pointSize, c)
}
