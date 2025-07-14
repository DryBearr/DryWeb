// ===============================================================
// File: drawing.go
// Description: Logic of drawing for the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"wasm/render"
)

func StartDrawingLoop() {
	go func() {
		var prev *render.Coordinate

		for {
			select {
			case c, ok := <-drawLineCoordinateChan:
				if !ok {
					return
				}

				if prev == nil {
					prev = &c
				} else {
					drawLine(alivePixel, *prev, c)
					prev = &c
				}

			case <-resetPrevPoint:
				prev = nil

			case point, ok := <-drawPointCoordinateChan:
				if ok {
					drawPoint(alivePixel, point)
				}
			}
		}
	}()
}

func AddLineCordinateQueue(c render.Coordinate) {
	drawLineCoordinateChan <- c
}

func AddPointCordinateQueue(c render.Coordinate) {
	drawPointCoordinateChan <- c
}

func ResetPrevPoint() {
	resetPrevPoint <- struct{}{}
}

func drawLine(pixel render.Pixel, start render.Coordinate, end render.Coordinate) {
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

	AddFrame(RenderFrame{Frame: &tempFrame, Size: tempSize, C: &render.Coordinate{X: minX, Y: minY}})
}

func drawPoint(pixel render.Pixel, c render.Coordinate) {
	setPixel(pixel, c)

	pointFrame := [][]render.Pixel{
		{
			pixel,
		},
	}

	AddFrame(RenderFrame{Frame: &pointFrame, Size: pointSize, C: &c})
}

func setBoard() {
	frameMutex.Lock()
	defer frameMutex.Unlock()

	frame2D = make([][]render.Pixel, size.Height)
	for row := range frame2D {
		frame2D[row] = make([]render.Pixel, size.Width)
		for column := range frame2D[row] {
			frame2D[row][column] = backgroundPixel
		}
	}

	AddFrame(RenderFrame{Frame: &frame2D, Size: size})
}

func setPixel(pixel render.Pixel, c render.Coordinate) {
	frameMutex.Lock()
	defer frameMutex.Unlock()

	frame2D[c.Y][c.X] = pixel
}

func getPixel(c render.Coordinate) render.Pixel {
	frameMutex.Lock()
	defer frameMutex.Unlock()

	pixel := frame2D[c.Y][c.X]

	return pixel
}
