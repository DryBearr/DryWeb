// ===============================================================
// File: core.go
// Description: Core of the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"sync"
	"time"
	"wasm/render"
)

type RenderFrame struct {
	Frame *[][]render.Pixel
	C     *render.Coordinate
	Size  render.Size
}

var (
	api render.Renderer

	//Renderer
	latency   time.Duration
	frameChan chan RenderFrame

	//Drawing vars
	size render.Size

	frameMutex sync.Mutex
	frame2D    [][]render.Pixel

	resetPrevPoint          chan struct{}
	drawPointCoordinateChan chan render.Coordinate
	drawLineCoordinateChan  chan render.Coordinate

	pointSize render.Size

	deadPixel       render.Pixel
	alivePixel      render.Pixel
	backgroundPixel render.Pixel

	//Population vars
	maxPopulationCordinate render.Coordinate
	aliveCells             chan map[render.Coordinate]any
)

func StartGame(renderer render.Renderer) error {
	api = renderer

	latency = 16 * time.Millisecond // 60 fps 1 / 60 ~= 16.6

	size = api.GetSize()

	drawLineCoordinateChan = make(chan render.Coordinate, 100)  //TODO: use passed param
	drawPointCoordinateChan = make(chan render.Coordinate, 100) //TODO: use passed param
	resetPrevPoint = make(chan struct{}, 1)
	frameChan = make(chan RenderFrame, 144) //TODO: use passed param

	pointSize = render.Size{Width: 1, Height: 1}

	deadPixel = render.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	alivePixel = render.Pixel{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	backgroundPixel = render.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	setBoard()

	api.RegisterResizeEventListener(changeSize)
	api.RegisterMouseDragEventListener(onDrag)
	api.RegisterMouseClickEventListener(onClick)
	api.RegisterMouseDragEndEventListener(onDragEnd)

	StartRenderLoop()
	StartDrawingLoop()

	return nil
}

func changeSize(s render.Size) error {
	if size.Width == s.Width && size.Height == s.Height {
		return nil
	}

	size = s
	setBoard()

	return nil
}

func onClick(c render.Coordinate) error {
	AddPointCordinateQueue(c)

	return nil
}

func onDrag(c render.Coordinate) error {
	AddLineCordinateQueue(c)

	return nil
}

func onDragEnd(c render.Coordinate) error {
	ResetPrevPoint()

	return nil
}
