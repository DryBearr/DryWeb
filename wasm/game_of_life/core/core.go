// ===============================================================
// File: core.go
// Description: Core of the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"sync"
	"time"
	"wasm/dryengine"
)

var (
	Engine dryengine.DryEngine

	//Drawing vars
	Size dryengine.Size

	FrameMutex sync.Mutex
	Frame2D    [][]dryengine.Pixel

	ResetPrevPointChan      chan struct{}
	DrawPointCoordinateChan chan dryengine.Coordinate2D
	DrawLineCoordinateChan  chan dryengine.Coordinate2D

	PointSize dryengine.Size

	DeadPixel       dryengine.Pixel
	AlivePixel      dryengine.Pixel
	BackgroundPixel dryengine.Pixel

	//Population vars
	BoundaryCordinate dryengine.Coordinate2D

	AliveCellsMutex sync.Mutex
	AliveCells      map[dryengine.Coordinate2D]any

	PopulationMutex  sync.Mutex
	PopulationCond   = sync.NewCond(&PopulationMutex)
	PausedPopulation = false

	//ResetPopulation chan struct{} //TODO: signal to clear the board
)

func StartGame(engine dryengine.DryEngine) {
	Engine = engine

	Size = Engine.GetSize()

	DrawLineCoordinateChan = make(chan dryengine.Coordinate2D, 100)  //TODO: use passed param
	DrawPointCoordinateChan = make(chan dryengine.Coordinate2D, 100) //TODO: use passed param
	ResetPrevPointChan = make(chan struct{}, 1)
	BoundaryCordinate = dryengine.Coordinate2D{X: 6000, Y: 6000}

	PointSize = dryengine.Size{Width: 1, Height: 1}

	DeadPixel = dryengine.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	AlivePixel = dryengine.Pixel{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	BackgroundPixel = dryengine.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	SetBoard()
	Engine.AddFrame(&dryengine.RenderFrame{Frame: &Frame2D, FrameSize: Size})

	Engine.RegisterResizeEventListener(ChangeSize)
	Engine.RegisterMouseDragEventListener(OnDrag)
	Engine.RegisterMouseClickEventListener(OnClick)
	Engine.RegisterMouseDragEndEventListener(OnDragEnd)

	Engine.StartRenderLoop(16 * time.Millisecond)

	RunPopulationLoop(100 * time.Millisecond)

	StartDrawingLoop()

	select {} //Run Forever when ever :3
}

func ChangeSize(s dryengine.Size) error {
	if Size.Width == s.Width && Size.Height == s.Height {
		return nil
	}

	Size = s
	SetBoard()

	return nil
}

func OnClick(c dryengine.Coordinate2D) error {
	AddPointCordinateQueue(c)
	return nil
}

func OnDrag(c dryengine.Coordinate2D) error {
	PausePopulation()
	AddLineCordinateQueue(c)
	return nil
}

func OnDragEnd(c dryengine.Coordinate2D) error {
	ResetPrevPoint()
	ResumePopulation()
	return nil
}
