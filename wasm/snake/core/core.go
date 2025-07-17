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
	"wasm/dryengine"
)

type Move dryengine.Coordinate2D

const (
	wall      byte = 4
	snakeHead byte = 3
	snakeTail byte = 2
	apple     byte = 1

	boardSize = 17 //original 15 but up down and left right + 2
	latency   = 16 //60 frame per second
)

var (
	engine dryengine.DryEngine

	boardMutex   sync.Mutex
	board        [][]byte
	pointsEarned int

	//TODO:
	pauseMu   sync.Mutex
	pause     bool
	pauseCond sync.Cond

	moveLeft  = Move{X: -1, Y: 0}
	moveUp    = Move{X: 0, Y: -1}
	moveDown  = Move{X: 0, Y: 1}
	moveRight = Move{X: 1, Y: 0}

	//TODO: i can just use board so future me fix this poop :)
	snakePartsMutex sync.Mutex
	delayedTail     *dryengine.Coordinate2D
	snakeParts      []dryengine.Coordinate2D

	appleMutex   sync.Mutex
	droppedApple *dryengine.Coordinate2D

	pointMutex sync.Mutex
	points     int

	endGameChan chan any

	snakeDirectionMutex sync.Mutex
	snakeDirection      Move = moveDown

	minimumDuration      = 100
	minusDuration        = 5
	maxDuration          = 200
	currentDuration      = maxDuration
	currentDurationMutex sync.Mutex
	ticker               = time.NewTicker(time.Millisecond * time.Duration(currentDuration))

	size      dryengine.Size
	sizeMutex sync.Mutex

	backgroundColor = dryengine.Pixel{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
	snakeColor = dryengine.Pixel{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	snackColor = dryengine.Pixel{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	}
	wallColor = dryengine.Pixel{
		R: 255,
		G: 0,
		B: 0,
		A: 50,
	}
)

func StartGame(newEngine dryengine.DryEngine) {
	initBoard()

	initSnake()

	droppedApple = &dryengine.Coordinate2D{X: rand.Intn(boardSize-2) + 1, Y: rand.Intn(boardSize-2) + 1}

	engine = newEngine

	size = engine.GetSize()

	endGameChan = make(chan any)

	engine.RegisterKeyDownEventListener(onKeyDown)
	engine.RegisterResizeEventListener(onResize)
	engine.RegisterSwipeEventListener(onSwipe)

	engine.StartRenderLoop(16 * time.Millisecond)

	startGameLoop()

	select {} //Run Forever when ever :3
}

//Getters & Setters with mutex

func getBoard() [][]byte {
	boardMutex.Lock()
	defer boardMutex.Unlock()

	return board
}

func setBoard(newBoard [][]byte) {
	boardMutex.Lock()
	defer boardMutex.Unlock()

	board = newBoard
}

func setSize(newSize dryengine.Size) {
	sizeMutex.Lock()
	defer sizeMutex.Unlock()

	size = newSize
}

func getSize() dryengine.Size {
	sizeMutex.Lock()
	defer sizeMutex.Unlock()

	return size
}

func setSnakeDirection(newDirection Move) {
	snakeDirectionMutex.Lock()
	defer snakeDirectionMutex.Unlock()

	snakeDirection = newDirection
}

func getSnakeDirection() Move {
	snakeDirectionMutex.Lock()
	defer snakeDirectionMutex.Unlock()

	return snakeDirection
}

func setSnakeParts(newSnakeParts []dryengine.Coordinate2D) {
	snakePartsMutex.Lock()
	defer snakePartsMutex.Unlock()

	snakeParts = newSnakeParts
}

func getSnakeParts() []dryengine.Coordinate2D {
	snakePartsMutex.Lock()
	defer snakePartsMutex.Unlock()

	return snakeParts
}

func getApple() *dryengine.Coordinate2D {
	appleMutex.Lock()
	defer appleMutex.Unlock()

	return droppedApple
}

func setApple(newApple *dryengine.Coordinate2D) {
	appleMutex.Lock()
	defer appleMutex.Unlock()

	droppedApple = newApple
}

//Init funcs for game

func initBoard() {
	newBoard := make([][]byte, boardSize)
	for row := range newBoard {
		newBoard[row] = make([]byte, boardSize)

		for column := range newBoard[row] {
			newBoard[row][column] = 0

			if row == boardSize-1 || row == 0 || column == boardSize-1 || column == 0 { //walls
				newBoard[row][column] = 4
			}
		}
	}

	newBoard[1][1] = snakeHead

	setBoard(newBoard)
}

func initSnake() {
	setSnakeParts([]dryengine.Coordinate2D{
		{
			X: 1,
			Y: 1,
		},
	})
}

// Game Logic funcs

func moveSnake(move Move) {
	currentSnakeParts := getSnakeParts()

	prevCoordinate := currentSnakeParts[0] //head
	currentSnakeParts[0].X += move.X
	currentSnakeParts[0].Y += move.Y

	for idx := range currentSnakeParts {
		if idx == 0 {
			continue
		}

		temp := currentSnakeParts[idx]
		currentSnakeParts[idx] = prevCoordinate
		prevCoordinate = temp

	}

	setSnakeParts(currentSnakeParts)
}

func checkState() {
	currentBoard := getBoard()
	currentSnakeParts := getSnakeParts()

	if delayedTail != nil {
		currentSnakeParts = append(currentSnakeParts, *delayedTail)
		delayedTail = nil

		setSnakeParts(currentSnakeParts)
	}

	switch board[snakeParts[0].Y][snakeParts[0].X] {
	case wall, snakeTail:
		go func() {
			endGameChan <- struct{}{}
		}()
		return
	case apple:
		decreaseDuration()

		newTail := currentSnakeParts[len(snakeParts)-1]
		delayedTail = &newTail

		setApple(&dryengine.Coordinate2D{X: rand.Intn(boardSize-2) + 1, Y: rand.Intn(boardSize-2) + 1})

		increasePoints()
	}

	for row := range currentBoard {
		for column := range currentBoard[row] {
			currentBoard[row][column] = 0
			if row == boardSize-1 || row == 0 || column == boardSize-1 || column == 0 { //walls
				currentBoard[row][column] = 4
			}
		}
	}

	currentDroppedApple := getApple()
	if getApple() != nil {
		currentBoard[currentDroppedApple.Y][currentDroppedApple.X] = apple
	}

	currentBoard[currentSnakeParts[0].Y][currentSnakeParts[0].X] = snakeHead

	for _, snakePart := range currentSnakeParts[1:] {
		currentBoard[snakePart.Y][snakePart.X] = snakeTail
	}

	setBoard(currentBoard)
}

func increasePoints() {
	pointMutex.Lock()
	defer pointMutex.Unlock()

	points += 1
}

func resetPoints() {
	pointMutex.Lock()
	defer pointMutex.Unlock()

	points = 0
}

func decreaseDuration() {
	currentDurationMutex.Lock()
	defer currentDurationMutex.Unlock()

	if currentDuration > minimumDuration {
		currentDuration -= minimumDuration
		ticker.Stop()
		ticker = time.NewTicker(time.Duration(currentDuration) * time.Millisecond)
	}
}

func resetDuration() {
	currentDurationMutex.Lock()
	defer currentDurationMutex.Unlock()

	currentDuration = maxDuration
	ticker.Stop()
	ticker = time.NewTicker(time.Duration(currentDuration) * time.Millisecond)
}

func startGameLoop() {
	go func() {
	GameLoop:
		for {
			select {
			case <-ticker.C:
				moveSnake(getSnakeDirection())

				checkState()

				engine.AddFrame(&dryengine.RenderFrame{
					Frame:     boardToFrame(),
					FrameSize: size,
				})
			case <-endGameChan:
				resetDuration()

				resetPoints()

				initBoard()

				initSnake()

				setSnakeDirection(moveDown)

				continue GameLoop
			}
		}
	}()
}

// Event handlers
func onSwipe(direction dryengine.SwipeDirection) error {
	currentSnakeDirection := getSnakeDirection()

	switch direction {
	case dryengine.SwipeUp:
		if currentSnakeDirection != moveDown {
			setSnakeDirection(moveUp)
		}
	case dryengine.SwipeLeft:
		if currentSnakeDirection != moveRight {
			setSnakeDirection(moveLeft)
		}
	case dryengine.SwipeRight:
		if currentSnakeDirection != moveLeft {
			setSnakeDirection(moveRight)
		}
	case dryengine.SwipeDown:
		if currentSnakeDirection != moveUp {
			setSnakeDirection(moveDown)
		}
	}

	return nil
}

func onKeyDown(key dryengine.Key) error {
	currentSnakeDirection := getSnakeDirection()

	switch key {
	case dryengine.WKey:
		if currentSnakeDirection != moveDown {
			setSnakeDirection(moveUp)
		}
	case dryengine.AKey:
		if currentSnakeDirection != moveRight {
			setSnakeDirection(moveLeft)
		}
	case dryengine.DKey:
		if currentSnakeDirection != moveLeft {
			setSnakeDirection(moveRight)
		}
	case dryengine.SKey:
		if currentSnakeDirection != moveUp {
			setSnakeDirection(moveDown)
		}
	}

	return nil
}

func onResize(newSize dryengine.Size) error {
	setSize(newSize)

	return nil
}

//Game rendering funcs

func boardToFrame() *[][]dryengine.Pixel {
	currentBoard := getBoard()

	currentSize := getSize()

	newFrame2D := make([][]dryengine.Pixel, currentSize.Height)
	for idx := range newFrame2D {
		newFrame2D[idx] = make([]dryengine.Pixel, currentSize.Width)
	}

	for frameYIndex := range newFrame2D {
		for frameXIndex := range newFrame2D[frameYIndex] {
			boardYIndex := frameYIndex * boardSize / currentSize.Height
			boardXIndex := frameXIndex * boardSize / currentSize.Width

			switch currentBoard[boardYIndex][boardXIndex] {
			case snakeTail, snakeHead:
				newFrame2D[frameYIndex][frameXIndex] = snakeColor
			case wall:
				newFrame2D[frameYIndex][frameXIndex] = wallColor
			case apple:
				newFrame2D[frameYIndex][frameXIndex] = snackColor
			default:
				newFrame2D[frameYIndex][frameXIndex] = backgroundColor
			}
		}
	}

	return &newFrame2D
}
