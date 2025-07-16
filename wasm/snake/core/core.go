// ===============================================================
// File: core.go
// Description: Snake game logic and rendering
// Author: DryBearr
// ===============================================================

package core

import (
	"math/rand"
	"sync"
	"time"
	"wasm/render"
)

type Move render.Coordinate

const (
	wall      byte = 4
	snakeHead byte = 3
	snakeTail byte = 2
	snack     byte = 1

	boardSize = 17 //original 15 but up down and left right + 2
	latency   = 16 //60 frame per second
)

var (
	api render.Renderer

	boardMutex   sync.Mutex
	board        [][]byte
	pointsEarned int

	pauseMu   sync.Mutex
	pause     bool
	pauseCond sync.Cond

	moveLeft  = Move{X: -1, Y: 0}
	moveUp    = Move{X: 0, Y: -1}
	moveDown  = Move{X: 0, Y: 1}
	moveRight = Move{X: 1, Y: 0}

	//TODO: i can just use board so future me fix this poop :)
	snakePartsMutex sync.Mutex
	delayedTail     *render.Coordinate
	snakeParts      []render.Coordinate

	snackMutex   sync.Mutex
	droppedSnack *render.Coordinate
	snackCount   int

	endGameChan chan any

	snakeDirectionMutex sync.Mutex
	snakeDirection      Move = moveDown

	minimumDuration      = 100
	minusDuration        = 5
	maxDuration          = 200
	currentDuration      = maxDuration
	currentDurationMutex sync.Mutex
	ticker               = time.NewTicker(time.Millisecond * time.Duration(currentDuration))

	size render.Size

	backgroundColor = render.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	snakeColor = render.Pixel{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	snackColor = render.Pixel{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	}

	wallColor = render.Pixel{
		R: 255,
		G: 0,
		B: 0,
		A: 50,
	}

	frameChan chan *[][]render.Pixel
)

func StartGame(renderer render.Renderer) {
	initBoard()
	initSnake()
	droppedSnack = &render.Coordinate{X: rand.Intn(boardSize-2) + 1, Y: rand.Intn(boardSize-2) + 1}

	api = renderer
	size = api.GetSize()

	endGameChan = make(chan any)
	frameChan = make(chan *[][]render.Pixel, 1000)

	api.RegisterKeyDownEventListener(onKeyDown)
	api.RegisterResizeEventListener(onResize)

	startRenderLoop()
	startGameLoop()

	select {}
}

func initBoard() {
	board = make([][]byte, boardSize)
	for row := range board {
		board[row] = make([]byte, boardSize)
		for column := range board[row] {
			board[row][column] = 0
			if row == boardSize-1 || row == 0 || column == boardSize-1 || column == 0 { //walls
				board[row][column] = 4
			}
		}
	}
}

func initSnake() {
	boardMutex.Lock()
	board[1][1] = snakeHead
	boardMutex.Unlock()

	snakePartsMutex.Lock()
	snakeParts = []render.Coordinate{
		{
			X: 1,
			Y: 1,
		},
	}
	snakePartsMutex.Unlock()
}

func moveSnake(move Move) {
	boardMutex.Lock()
	snakePartsMutex.Lock()

	prevCoordinate := snakeParts[0] //head
	snakeParts[0].X += move.X
	snakeParts[0].Y += move.Y

	for idx := range snakeParts {
		if idx == 0 {
			continue
		}

		temp := snakeParts[idx]
		snakeParts[idx] = prevCoordinate
		prevCoordinate = temp

	}
	snakePartsMutex.Unlock()
	boardMutex.Unlock()
}

func checkState() {
	boardMutex.Lock()
	snakePartsMutex.Lock()

	if delayedTail != nil {
		snakeParts = append(snakeParts, *delayedTail)
		delayedTail = nil
	}

	switch board[snakeParts[0].Y][snakeParts[0].X] {
	case wall, snakeTail:
		go func() {
			endGameChan <- struct{}{}
		}()

		snakePartsMutex.Unlock()
		boardMutex.Unlock()
		return
	case snack:
		decreaseDuration()

		newTail := snakeParts[len(snakeParts)-1]
		delayedTail = &newTail

		snackMutex.Lock()
		droppedSnack = &render.Coordinate{X: rand.Intn(boardSize-2) + 1, Y: rand.Intn(boardSize-2) + 1}
		snackCount += 1
		snackMutex.Unlock()
	}

	for row := range board {
		for column := range board[row] {
			board[row][column] = 0
			if row == boardSize-1 || row == 0 || column == boardSize-1 || column == 0 { //walls
				board[row][column] = 4
			}
		}
	}

	snackMutex.Lock()
	if droppedSnack != nil {
		board[droppedSnack.Y][droppedSnack.X] = snack
	}
	snackMutex.Unlock()

	board[snakeParts[0].Y][snakeParts[0].X] = snakeHead

	for _, snakePart := range snakeParts[1:] {
		board[snakePart.Y][snakePart.X] = snakeTail
	}

	snakePartsMutex.Unlock()
	boardMutex.Unlock()
}

func startGameLoop() {
	go func() {
	GameLoop:
		for {
			select {
			case <-ticker.C:
				moveSnake(snakeDirection)
				checkState()
				go func() {
					frameChan <- BoardToFrame()
				}()
			case <-endGameChan:
				currentDurationMutex.Lock()
				currentDuration = maxDuration
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(currentDuration) * time.Millisecond)
				currentDurationMutex.Unlock()

				initBoard()
				initSnake()
				snakeDirectionMutex.Lock()
				snakeDirection = moveDown
				snakeDirectionMutex.Unlock()
				continue GameLoop
			}
		}
	}()
}

func changeSnakeDirection(move Move) {
	snakeDirectionMutex.Lock()
	defer snakeDirectionMutex.Unlock()

	snakeDirection = move
}

func onKeyDown(key render.Key) error {
	switch key {
	case render.WKey:
		if snakeDirection != moveDown {
			changeSnakeDirection(moveUp)
		}
	case render.AKey:
		if snakeDirection != moveRight {
			changeSnakeDirection(moveLeft)
		}
	case render.DKey:
		if snakeDirection != moveLeft {
			changeSnakeDirection(moveRight)
		}
	case render.SKey:
		if snakeDirection != moveUp {
			changeSnakeDirection(moveDown)
		}
	}

	return nil
}

func decreaseDuration() {
	currentDurationMutex.Lock()
	if currentDuration > minimumDuration {
		currentDuration -= minimumDuration
		ticker.Stop()
		ticker = time.NewTicker(time.Duration(currentDuration) * time.Millisecond)
	}
	currentDurationMutex.Unlock()
}

func onResize(newSize render.Size) error {
	return nil
}

func BoardToFrame() *[][]render.Pixel {
	boardMutex.Lock()
	defer boardMutex.Unlock()

	newFrame2D := make([][]render.Pixel, size.Height)
	for idx := range newFrame2D {
		newFrame2D[idx] = make([]render.Pixel, size.Width)
	}

	for frameYIndex := range newFrame2D {
		for frameXIndex := range newFrame2D[frameYIndex] {
			boardYIndex := frameYIndex * boardSize / size.Height
			boardXIndex := frameXIndex * boardSize / size.Width
			switch board[boardYIndex][boardXIndex] {
			case snakeTail, snakeHead:
				newFrame2D[frameYIndex][frameXIndex] = snakeColor
			case wall:
				newFrame2D[frameYIndex][frameXIndex] = wallColor
			case snack:
				newFrame2D[frameYIndex][frameXIndex] = snackColor
			default:
				newFrame2D[frameYIndex][frameXIndex] = backgroundColor
			}
		}
	}

	return &newFrame2D
}

func startRenderLoop() {
	go func() {
		timer := time.NewTimer(latency)
		defer timer.Stop()

		for {
			select {
			case frame, ok := <-frameChan:
				if !ok {
					return
				}

				api.DrawFrame(frame, render.Size{
					Height: len(*frame),
					Width:  len((*frame)[0]),
				})

				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(latency)

			case <-timer.C:
				timer.Reset(latency)
			}
		}
	}()
}
