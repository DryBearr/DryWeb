// ===============================================================
// File: population.go
// Description: Logic for tracking population in the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"time"
	"wasm/render"
)

func RunPopulationLoop(populateInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(populateInterval)
		defer ticker.Stop()

		for {
			ticker.Reset(populateInterval)

			PopulationMutex.Lock()
			for PausedPopulation {
				PopulationCond.Wait()
			}
			PopulationMutex.Unlock()

			select {
			case <-ticker.C:
				PopulateFrame()
				//TODO:
				//case <-reset:
			}
		}
	}()
}

func PopulateFrame() {
	AliveCellsMutex.Lock()
	defer AliveCellsMutex.Unlock()

	possibleAliveCells := make(map[render.Coordinate]int)

	newAliveCells := make(map[render.Coordinate]any)

	for aliveCell := range AliveCells {
		nCoordinates := getNeighbourCoordinates(aliveCell, BoundaryCordinate.X, BoundaryCordinate.Y)

		aliveCellCount := 0

		for _, nCoordinate := range nCoordinates {
			if _, ok := AliveCells[nCoordinate]; ok {
				aliveCellCount += 1
			} else {
				possibleAliveCells[nCoordinate]++
			}
		}

		if aliveCellCount == 2 || aliveCellCount == 3 {
			newAliveCells[aliveCell] = struct{}{}
		}
	}

	for possibleAliveCell, aliveNeighbourCount := range possibleAliveCells {
		if aliveNeighbourCount == 3 {
			newAliveCells[possibleAliveCell] = struct{}{}
		}
	}

	FrameMutex.Lock()

	AliveCells = newAliveCells

	//TODO: frame and grid of living cells are not the same size so create translator for that
	tempCoordinate := render.Coordinate{}
	for y := range Frame2D {
		for x := range Frame2D[y] {
			tempCoordinate.Y = y
			tempCoordinate.X = x

			if _, ok := AliveCells[tempCoordinate]; ok {
				Frame2D[y][x] = AlivePixel
			} else {
				Frame2D[y][x] = DeadPixel
			}
		}
	}

	//TODO: this is poop code
	AddFrame(RenderFrame{
		Frame: &Frame2D,
		Size:  Size,
	})

	FrameMutex.Unlock()
}

func ResurectCell(c render.Coordinate) {
	AliveCellsMutex.Lock()
	defer AliveCellsMutex.Unlock()

	if c.X < BoundaryCordinate.X && c.Y < BoundaryCordinate.Y {
		AliveCells[c] = struct{}{}
	}
}

func ResurectCellMany(coordinates []render.Coordinate) {
	AliveCellsMutex.Lock()
	defer AliveCellsMutex.Unlock()

	for _, c := range coordinates {
		if c.X < BoundaryCordinate.X && c.Y < BoundaryCordinate.Y {
			AliveCells[c] = struct{}{}
		}
	}
}

func getNeighbourCoordinates(c render.Coordinate, width, height int) []render.Coordinate {
	neighbors := make([]render.Coordinate, 0, 8)

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx := c.X + dx
			ny := c.Y + dy

			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				neighbors = append(neighbors, render.Coordinate{
					X: nx,
					Y: ny,
				})
			}
		}
	}

	return neighbors
}

func PausePopulation() {
	PopulationMutex.Lock()
	defer PopulationMutex.Unlock()

	PausedPopulation = true
}

func ResumePopulation() {
	PopulationMutex.Lock()
	defer PopulationMutex.Unlock()

	PausedPopulation = false

	PopulationCond.Broadcast()
}
