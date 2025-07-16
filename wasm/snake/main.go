// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
	"wasm/render"
	"wasm/webrender"
)

func main() {
	api := render.Renderer(webrender.Api)
	size := api.GetSize()

	var frameMutex sync.Mutex

	frame2D := make([][]render.Pixel, size.Height)
	for idx := range frame2D {
		frame2D[idx] = make([]render.Pixel, size.Width)
	}

	changeSize := func(s render.Size) error {
		frameMutex.Lock()
		defer frameMutex.Unlock()

		if size.Width == s.Width && size.Height == s.Height {
			return nil
		}
		size = s

		frame2D = make([][]render.Pixel, s.Height)
		for idx := range frame2D {
			frame2D[idx] = make([]render.Pixel, s.Width)
		}
		return nil
	}

	keyDown := func(key render.Key) error {
		log.Println(key)
		return nil
	}

	api.RegisterResizeEventListener(changeSize)
	api.RegisterKeyDownEventListener(keyDown)

	for {
		SetRandomFrame(&frame2D)

		time.Sleep(16 * time.Millisecond)

		api.DrawFrame(&frame2D, size)
	}
}

// Any live cell with fewer than two live neighbours dies (referred to as underpopulation).
//
// Any live cell with more than three live neighbours dies (referred to as overpopulation).
// Any live cell with two or three live neighbours lives, unchanged, to the next generation.
// Any dead cell with exactly three live neighbours comes to life.

func SetRandomFrame(frame2D *[][]render.Pixel) {
	for _, frame1D := range *frame2D {
		for idx, pixel := range frame1D {
			if rand.Intn(3)%2 == 0 {
				pixel.R = 255
				pixel.G = 255
				pixel.B = 255
			} else {
				pixel.R = 0
				pixel.G = 0
				pixel.B = 0
			}
			pixel.A = 255

			frame1D[idx] = pixel
		}
	}
}
