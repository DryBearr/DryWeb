// ===============================================================
// File: core.go
// Description: Core of the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"sync"
	"wasm/render"
)

var (
	api  render.Renderer
	size render.Size

	frameMutex sync.Mutex
	frame2D    [][]render.Pixel
)

func StartGame(renderer render.Renderer) error {
	api = renderer
	size = api.GetSize()
	frame2D = make([][]render.Pixel, size.Height)
	for idx := range frame2D {
		frame2D[idx] = make([]render.Pixel, size.Width)
	}
	SetBlackBoard()

	api.DrawFrame(&frame2D, size)

	changeSize := func(s render.Size) error {
		if size.Width == s.Width && size.Height == s.Height {
			return nil
		}

		size = s

		frameMutex.Lock()
		frame2D = make([][]render.Pixel, s.Height)
		for idx := range frame2D {
			frame2D[idx] = make([]render.Pixel, s.Width)
		}
		frameMutex.Unlock()

		SetBlackBoard()

		api.DrawFrame(&frame2D, size)
		return nil
	}

	onDrag := func(c render.Coordinate) error {
		AddQueue(c)

		return nil
	}

	api.RegisterResizeEventListener(changeSize)
	api.RegisterMouseDragEventListener(onDrag)

	StartDrawingLoop()

	return nil
}

func SetBlackBoard() {
	frameMutex.Lock()
	defer frameMutex.Unlock()

	for _, frame1D := range frame2D {
		for idx, pixel := range frame1D {
			pixel.R = 0
			pixel.G = 0
			pixel.B = 0
			pixel.A = 255
			(frame1D)[idx] = pixel
		}
	}
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
