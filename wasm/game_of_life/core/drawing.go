// ===============================================================
// File: drawing.go
// Description: Logic of drawing for the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"wasm/dryengine"
)

func StartDrawingLoop() {
	go func() {
		var prev *dryengine.Coordinate2D

		for {
			select {
			case c, ok := <-DrawLineCoordinateChan:
				if !ok {
					return
				}

				if prev == nil {
					prev = &c
				} else {
					predictedCoordinates := DrawLine(AlivePixel, *prev, c)

					ResurectCellMany(predictedCoordinates) //TODO: this mf is a black sheep so move somewhere else

					prev = &c
				}

			case <-ResetPrevPointChan:
				prev = nil

			case point, ok := <-DrawPointCoordinateChan:
				if ok {
					DrawPoint(AlivePixel, point)
					ResurectCell(point) //TODO: this mf is a black sheep so move somewhere else
				}
			}
		}
	}()
}

func AddLineCordinateQueue(c dryengine.Coordinate2D) {
	DrawLineCoordinateChan <- c
}

func AddPointCordinateQueue(c dryengine.Coordinate2D) {
	DrawPointCoordinateChan <- c
}

func ResetPrevPoint() {
	ResetPrevPointChan <- struct{}{}
}

func DrawLine(pixel dryengine.Pixel, start dryengine.Coordinate2D, end dryengine.Coordinate2D) []dryengine.Coordinate2D {
	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y

	diffX := Abs(x0 - x1)
	diffY := Abs(y0 - y1)

	minX := min(x0, x1)
	minY := min(y0, y1)

	tempSize := dryengine.Size{
		Width:  diffX + 1,
		Height: diffY + 1,
	}

	reserveSize := max(diffX, diffY)

	ultraInstinctCoordinates := make([]dryengine.Coordinate2D, 0, reserveSize) //predicted coordinates between start and end points

	tempFrame := make([][]dryengine.Pixel, tempSize.Height)

	FrameMutex.Lock()
	for row := range tempFrame {
		tempFrame[row] = make([]dryengine.Pixel, tempSize.Width)

		for column := range tempFrame[row] {
			frameY := minY + row
			frameX := minX + column

			if frameY >= 0 && frameY < len(Frame2D) && frameX >= 0 && frameX < len(Frame2D[0]) {
				tempFrame[row][column] = Frame2D[frameY][frameX]
			}
		}
	}
	FrameMutex.Unlock()

	stepX := 1
	if x0 > x1 {
		stepX = -1
	}

	stepY := 1
	if y0 > y1 {
		stepY = -1
	}

	err := diffX - diffY

	FrameMutex.Lock()
	for {
		// Compute tempFrame indices
		tempX := x0 - minX
		tempY := y0 - minY

		// Check tempFrame bounds
		if tempY >= 0 && tempY < tempSize.Height && tempX >= 0 && tempX < tempSize.Width {
			tempFrame[tempY][tempX] = pixel
		}

		// Update main frame if needed
		if y0 >= 0 && y0 < len(Frame2D) && x0 >= 0 && x0 < len(Frame2D[0]) {
			Frame2D[y0][x0] = pixel
			ultraInstinctCoordinates = append(ultraInstinctCoordinates, dryengine.Coordinate2D{X: x0, Y: y0})
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
	FrameMutex.Unlock()
	Engine.AddFrame(&dryengine.RenderFrame{Frame: &tempFrame, FrameSize: tempSize, C: &dryengine.Coordinate2D{X: minX, Y: minY}})

	return ultraInstinctCoordinates
}

func DrawPoint(pixel dryengine.Pixel, c dryengine.Coordinate2D) {
	SetPixel(pixel, c)

	pointFrame := [][]dryengine.Pixel{
		{
			pixel,
		},
	}

	Engine.AddFrame(&dryengine.RenderFrame{Frame: &pointFrame, FrameSize: PointSize, C: &c})
}

func SetBoard() {
	FrameMutex.Lock()
	defer FrameMutex.Unlock()

	Frame2D = make([][]dryengine.Pixel, Size.Height)
	for row := range Frame2D {
		Frame2D[row] = make([]dryengine.Pixel, Size.Width)
		for column := range Frame2D[row] {
			Frame2D[row][column] = BackgroundPixel
		}
	}
}

func SetPixel(pixel dryengine.Pixel, c dryengine.Coordinate2D) {
	FrameMutex.Lock()
	defer FrameMutex.Unlock()

	Frame2D[c.Y][c.X] = pixel
}
