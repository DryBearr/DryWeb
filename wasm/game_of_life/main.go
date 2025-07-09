// ===============================================================
// File: main.go
// Description: application's entry point
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package main

import (
	"math/rand"
	"sync"
	"time"
	"wasm/render"
)

func main() {
	size := render.Api.GetSize()

	var frameMutex sync.Mutex

	frame1D := make([]render.Pixel, size.Width*size.Height)

	changeSize := func(s render.Size) error {
		frameMutex.Lock()
		defer frameMutex.Unlock()

		if size.Width == s.Width && size.Height == s.Height {
			return nil
		}
		size = s

		frame1D = make([]render.Pixel, s.Width*s.Height)

		return nil
	}

	render.Api.RegisterResizeEventListener(changeSize)

	SetRandomFrame(&frame1D)

	for {
		SetRandomFrame(&frame1D)

		time.Sleep(16 * time.Millisecond)

		render.Api.DrawFrame(&frame1D, size)
	}
}

func SetRandomFrame(frame1D *[]render.Pixel) {
	for idx, pixel := range *frame1D {
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

		(*frame1D)[idx] = pixel
	}
}
