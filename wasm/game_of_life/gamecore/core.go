// ===============================================================
// File: core.go
// Description: Core of the game of life
// Author: DryBearr
// ===============================================================

package gamecore

import (
	"sync"
	"time"
	"wasm/dryeve/engine"
	"wasm/dryeve/models"
)

var (
	Engine engine.Engine

	//Drawing vars
	Width  int = 800
	Height int = 600

	FrameMutex sync.Mutex
	Frame2D    [][]models.Pixel

	ResetPrevPointChan      chan struct{}
	DrawPointCoordinateChan chan models.Point2D
	DrawLineCoordinateChan  chan models.Point2D

	DeadPixel       models.Pixel
	AlivePixel      models.Pixel
	BackgroundPixel models.Pixel

	//Population vars
	BoundaryCordinate models.Point2D

	AliveCellsMutex sync.Mutex
	AliveCells      map[models.Point2D]any

	PopulationMutex  sync.Mutex
	PopulationCond   = sync.NewCond(&PopulationMutex)
	PausedPopulation = false

	//ResetPopulation chan struct{} //TODO: signal to clear the board
)

func StartGame(engine engine.Engine) {
	Engine = engine

	DrawLineCoordinateChan = make(chan models.Point2D, 100)  //TODO: use passed param
	DrawPointCoordinateChan = make(chan models.Point2D, 100) //TODO: use passed param
	ResetPrevPointChan = make(chan struct{}, 1)
	BoundaryCordinate = models.Point2D{X: 6000, Y: 6000}

	DeadPixel = models.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	AlivePixel = models.Pixel{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	BackgroundPixel = models.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	SetBoard()
	Engine.AddFrame(models.RenderFrame{Frame: &Frame2D})

	Engine.Events.RegisterResizeEventListener(ChangeSize)
	Engine.Events.RegisterMouseDragEventListener(OnDrag)
	Engine.Events.RegisterMouseClickEventListener(OnClick)
	Engine.Events.RegisterMouseDragEndEventListener(OnDragEnd)

	Engine.StartRenderLoop()

	RunPopulationLoop(100 * time.Millisecond)

	StartDrawingLoop()

	select {} //Run Forever when ever :3
}

func ChangeSize(newWidth int, newHeight int) error {
	if Width == newWidth && Height == newHeight {
		return nil
	}

	Width = newWidth
	Height = newHeight

	SetBoard()

	return nil
}

func OnClick(c models.Point2D) error {
	AddPointCordinateQueue(c)
	return nil
}

func OnDrag(c models.Point2D) error {
	PausePopulation()
	AddLineCordinateQueue(c)
	return nil
}

func OnDragEnd(c models.Point2D) error {
	ResetPrevPoint()
	ResumePopulation()
	return nil
}
