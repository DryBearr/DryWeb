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
	"wasm/render"
)

func main() {
	width := render.GetWidth()
	height := render.GetHeight()

	frame1D := make([]render.Pixel, width*height)

	SetRandomFrame(&frame1D)
	for {
		log.Println("yee")
		SetRandomFrame(&frame1D)
		render.AddFrame(&frame1D)
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
