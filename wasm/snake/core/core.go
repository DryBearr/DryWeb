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
	apple     byte = 1

	boardSize = 17 //original 15 but up down and left right + 2
	latency   = 16 //60 frame per second
)

var (
	api render.Renderer

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
	delayedTail     *render.Coordinate
	snakeParts      []render.Coordinate

	appleMutex   sync.Mutex
	droppedApple *render.Coordinate

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

	size      render.Size
	sizeMutex sync.Mutex

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
	droppedApple = &render.Coordinate{X: rand.Intn(boardSize-2) + 1, Y: rand.Intn(boardSize-2) + 1}

	api = renderer
	size = api.GetSize()

	endGameChan = make(chan any)
	frameChan = make(chan *[][]render.Pixel, 1000)

	api.RegisterKeyDownEventListener(onKeyDown)
	api.RegisterResizeEventListener(onResize)
	api.RegisterSwipeEventListener(onSwipe)

	startRenderLoop()
	startGameLoop()

	select {}
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

func setSize(newSize render.Size) {
	sizeMutex.Lock()
	defer sizeMutex.Unlock()

	size = newSize
}

func getSize() render.Size {
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

func setSnakeParts(newSnakeParts []render.Coordinate) {
	snakePartsMutex.Lock()
	defer snakePartsMutex.Unlock()

	snakeParts = newSnakeParts
}

func getSnakeParts() []render.Coordinate {
	snakePartsMutex.Lock()
	defer snakePartsMutex.Unlock()

	return snakeParts
}

func getApple() *render.Coordinate {
	appleMutex.Lock()
	defer appleMutex.Unlock()

	return droppedApple
}

func setApple(newApple *render.Coordinate) {
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
	setSnakeParts([]render.Coordinate{
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

		setApple(&render.Coordinate{X: rand.Intn(boardSize-2) + 1, Y: rand.Intn(boardSize-2) + 1})

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

				go func() {
					frameChan <- boardToFrame()
				}()
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
func onSwipe(direction render.SwipeDirection) error {
	currentSnakeDirection := getSnakeDirection()

	switch direction {
	case render.SwipeUp:
		if currentSnakeDirection != moveDown {
			setSnakeDirection(moveUp)
		}
	case render.SwipeLeft:
		if currentSnakeDirection != moveRight {
			setSnakeDirection(moveLeft)
		}
	case render.SwipeRight:
		if currentSnakeDirection != moveLeft {
			setSnakeDirection(moveRight)
		}
	case render.SwipeDown:
		if currentSnakeDirection != moveUp {
			setSnakeDirection(moveDown)
		}
	}

	return nil
}

func onKeyDown(key render.Key) error {
	currentSnakeDirection := getSnakeDirection()

	switch key {
	case render.WKey:
		if currentSnakeDirection != moveDown {
			setSnakeDirection(moveUp)
		}
	case render.AKey:
		if currentSnakeDirection != moveRight {
			setSnakeDirection(moveLeft)
		}
	case render.DKey:
		if currentSnakeDirection != moveLeft {
			setSnakeDirection(moveRight)
		}
	case render.SKey:
		if currentSnakeDirection != moveUp {
			setSnakeDirection(moveDown)
		}
	}

	return nil
}

func onResize(newSize render.Size) error {
	setSize(newSize)

	return nil
}

//Game rendering funcs

func boardToFrame() *[][]render.Pixel {
	currentBoard := getBoard()

	currentSize := getSize()

	newFrame2D := make([][]render.Pixel, currentSize.Height)
	for idx := range newFrame2D {
		newFrame2D[idx] = make([]render.Pixel, currentSize.Width)
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
