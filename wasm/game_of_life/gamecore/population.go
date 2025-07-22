// ===============================================================
// File: population.go
// Description: Logic for tracking population in the game of life
// Author: DryBearr
// ===============================================================

package gamecore

import (
	"time"
	"wasm/dryeve/models"
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

	possibleAliveCells := make(map[models.Point2D]int)

	newAliveCells := make(map[models.Point2D]any)

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
	tempCoordinate := models.Point2D{}
	for y := range Frame2D {
		for x := range Frame2D[y] {
			tempCoordinate.Y = float32(y)
			tempCoordinate.X = float32(x)

			if _, ok := AliveCells[tempCoordinate]; ok {
				Frame2D[y][x] = AlivePixel
			} else {
				Frame2D[y][x] = DeadPixel
			}
		}
	}

	//TODO: this is poop code
	Engine.AddFrame(models.RenderFrame{
		Frame: &Frame2D,
	})

	FrameMutex.Unlock()
}

func ResurectCell(c models.Point2D) {
	AliveCellsMutex.Lock()
	defer AliveCellsMutex.Unlock()

	if c.X < BoundaryCordinate.X && c.Y < BoundaryCordinate.Y {
		AliveCells[c] = struct{}{}
	}

}

func ResurectCellMany(coordinates []models.Point2D) {
	AliveCellsMutex.Lock()
	defer AliveCellsMutex.Unlock()

	for _, c := range coordinates {
		if c.X < BoundaryCordinate.X && c.Y < BoundaryCordinate.Y {
			AliveCells[c] = struct{}{}
		}
	}
}

func getNeighbourCoordinates(c models.Point2D, width, height float32) []models.Point2D {
	neighbors := make([]models.Point2D, 0, 8)

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx := c.X + float32(dx)
			ny := c.Y + float32(dy)

			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				neighbors = append(neighbors, models.Point2D{
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
