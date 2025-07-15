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
	Api render.Renderer

	//Renderer
	Latency   time.Duration
	FrameChan chan RenderFrame

	//Drawing vars
	Size render.Size

	FrameMutex sync.Mutex
	Frame2D    [][]render.Pixel

	ResetPrevPointChan      chan struct{}
	DrawPointCoordinateChan chan render.Coordinate
	DrawLineCoordinateChan  chan render.Coordinate

	PointSize render.Size

	DeadPixel       render.Pixel
	AlivePixel      render.Pixel
	BackgroundPixel render.Pixel

	//Population vars
	BoundaryCordinate render.Coordinate

	AliveCellsMutex sync.Mutex
	AliveCells      map[render.Coordinate]any

	PopulationMutex  sync.Mutex
	PopulationCond   = sync.NewCond(&PopulationMutex)
	PausedPopulation = false

	//ResetPopulation chan struct{} //TODO: signal to clear the board
)

func StartGame(renderer render.Renderer) {
	Api = renderer

	Latency = 16 * time.Millisecond // 60 fps 1 / 60 ~= 16.6

	Size = Api.GetSize()

	DrawLineCoordinateChan = make(chan render.Coordinate, 100)  //TODO: use passed param
	DrawPointCoordinateChan = make(chan render.Coordinate, 100) //TODO: use passed param
	ResetPrevPointChan = make(chan struct{}, 1)
	BoundaryCordinate = render.Coordinate{X: 6000, Y: 6000}
	FrameChan = make(chan RenderFrame, 144) //TODO: use passed param

	PointSize = render.Size{Width: 1, Height: 1}

	DeadPixel = render.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	AlivePixel = render.Pixel{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	BackgroundPixel = render.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	SetBoard()
	AddFrame(RenderFrame{Frame: &Frame2D, Size: Size})

	Api.RegisterResizeEventListener(ChangeSize)
	Api.RegisterMouseDragEventListener(OnDrag)
	Api.RegisterMouseClickEventListener(OnClick)
	Api.RegisterMouseDragEndEventListener(OnDragEnd)

	StartRenderLoop()
	RunPopulationLoop(100 * time.Millisecond)
	StartDrawingLoop()

	select {} //Run Forever when ever :3
}

func ChangeSize(s render.Size) error {
	if Size.Width == s.Width && Size.Height == s.Height {
		return nil
	}

	Size = s
	SetBoard()

	return nil
}

func OnClick(c render.Coordinate) error {
	AddPointCordinateQueue(c)
	return nil
}

func OnDrag(c render.Coordinate) error {
	PausePopulation()
	AddLineCordinateQueue(c)
	return nil
}

func OnDragEnd(c render.Coordinate) error {
	ResetPrevPoint()
	ResumePopulation()
	return nil
}
