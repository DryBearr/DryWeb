// ===============================================================
// File: drawing.go
// Description: Logic of drawing for the game of life
// Author: DryBearr
// ===============================================================

package gamecore

import (
	"wasm/dryeve/models"
)

func StartDrawingLoop() {
	go func() {
		var prev *models.Point2D

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

func AddLineCordinateQueue(c models.Point2D) {
	DrawLineCoordinateChan <- c
}

func AddPointCordinateQueue(c models.Point2D) {
	DrawPointCoordinateChan <- c
}

func ResetPrevPoint() {
	ResetPrevPointChan <- struct{}{}
}

func DrawLine(pixel models.Pixel, start models.Point2D, end models.Point2D) []models.Point2D {
	x0, y0 := int(start.X), int(start.Y)
	x1, y1 := int(end.X), int(end.Y)

	diffX := Abs(x0 - x1)
	diffY := Abs(y0 - y1)

	minX := min(x0, x1)
	minY := min(y0, y1)

	tempWidth := diffX + 1
	tempHeight := diffY + 1

	reserveSize := max(diffX, diffY)

	ultraInstinctCoordinates := make([]models.Point2D, 0, reserveSize) //predicted coordinates between start and end points

	tempFrame := make([][]models.Pixel, tempHeight)

	FrameMutex.Lock()
	for row := range tempFrame {
		tempFrame[row] = make([]models.Pixel, tempWidth)

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
		if tempY >= 0 && tempY < tempHeight && tempX >= 0 && tempX < tempWidth {
			tempFrame[tempY][tempX] = pixel
		}

		// Update main frame if needed
		if y0 >= 0 && y0 < len(Frame2D) && x0 >= 0 && x0 < len(Frame2D[0]) {
			Frame2D[y0][x0] = pixel
			ultraInstinctCoordinates = append(ultraInstinctCoordinates, models.Point2D{X: float32(x0), Y: float32(y0)})
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
	Engine.AddFrame(models.RenderFrame{Frame: &tempFrame, C: &models.Point2D{X: float32(minX), Y: float32(minY)}})

	return ultraInstinctCoordinates
}

func DrawPoint(pixel models.Pixel, c models.Point2D) {
	SetPixel(pixel, c)

	pointFrame := [][]models.Pixel{
		{
			pixel,
		},
	}

	Engine.AddFrame(models.RenderFrame{Frame: &pointFrame, C: &c})
}

func SetBoard() {
	FrameMutex.Lock()
	defer FrameMutex.Unlock()

	Frame2D = make([][]models.Pixel, Height)
	for row := range Frame2D {
		Frame2D[row] = make([]models.Pixel, Width)
		for column := range Frame2D[row] {
			Frame2D[row][column] = BackgroundPixel
		}
	}
}

func SetPixel(pixel models.Pixel, c models.Point2D) {
	FrameMutex.Lock()
	defer FrameMutex.Unlock()

	Frame2D[int(c.Y)][int(c.X)] = pixel
}
