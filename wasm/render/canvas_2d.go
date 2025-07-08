// ===============================================================
// File: canvas_2d.go
// Description: Functions for rendering in html canvas in 2d
// Author: DryBearr
// ===============================================================

//go:build js && wasm

// TODO: package godoc
package render

import (
	"context"
	"fmt"
	"log"
	"sync"
	"syscall/js"
)

// Models for rendering
type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Variables
var (
	height int
	width  int

	pixelBuf []byte

	frameQueue1D      []*[]Pixel
	frameQueue1DMutex sync.Mutex

	animationLoopCtx context.Context
	CancelAnimation  context.CancelFunc

	running      bool
	runningMutex sync.Mutex
)

// Init variables
func init() {
	//TODO: runtime changeable width and height
	if isWorker() {
		params := js.Global().Get("computeParams")
		if !params.Truthy() {
			panic("computeParams is not defined in the worker")
		}

		width = params.Get("width").Int()

		height = params.Get("height").Int()
	} else {
		panic("can't run outside worker!!!")
	}

	numPixels := width * height * 4

	pixelBuf = make([]byte, numPixels)

	running = false
}

func GetWidth() int {
	return width
}

func GetHeight() int {
	return height
}

func AddFrame(frame1D *[]Pixel) {
	addFrameQueue1D(frame1D)

	startAnimationLoop()
}

func sendFrameToWorker(frame1D *[]Pixel, numRoutines int) error {
	log.Println("[canvas_2d.go] sending frame to worker.")
	frameQueue1DMutex.Lock()
	defer frameQueue1DMutex.Unlock()

	if frame1D == nil {
		return fmt.Errorf("frame1D is nil")
	}

	total := len(*frame1D)
	if total == 0 {
		return fmt.Errorf("frame1D has no pixels")
	} else if total != height*width {
		return fmt.Errorf("frame1D must have same number of pixels as (width * height) has")
	}

	if numRoutines <= 0 {
		return fmt.Errorf("numWorkers must be > 0")
	}

	if numRoutines > total {
		numRoutines = total
	}

	var wg sync.WaitGroup

	chunkSize := total / numRoutines

	for worker := range numRoutines {
		start := worker * chunkSize

		end := start + chunkSize
		if worker == numRoutines-1 {
			end = total // last worker does any leftover
		}

		wg.Add(1)
		go func(start int, end int) {
			defer wg.Done()

			for idx := start; idx < end; idx++ {
				i := idx * 4
				pixelBuf[i+0] = (*frame1D)[idx].R
				pixelBuf[i+1] = (*frame1D)[idx].G
				pixelBuf[i+2] = (*frame1D)[idx].B
				pixelBuf[i+3] = (*frame1D)[idx].A
			}
		}(start, end)
	}

	wg.Wait()

	jsBuf := js.Global().Get("Uint8Array").New(len(pixelBuf))

	js.CopyBytesToJS(jsBuf, pixelBuf)

	js.Global().Call("postMessage", map[string]any{
		"type":   "pixels",
		"width":  width,
		"height": height,
		"pixels": jsBuf,
	})

	log.Println("[canvas_2d.go] sended message to worker success.")
	return nil
}

func startAnimationLoop() {
	if getRunning() {

		return
	}

	animationLoopCtx, CancelAnimation = context.WithCancel(context.Background())

	setRunning(true)

	go func() {
		for {
			select {
			case <-animationLoopCtx.Done():
				setRunning(false)

				return
			default:
			}

			if isEmptyQueue() {
				setRunning(false)

				return
			}

			frame := getFrameQueue1D()

			err := sendFrameToWorker(frame, 4) //TODO: set routines number base on how heavy is work to do
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func setRunning(state bool) {
	runningMutex.Lock()

	running = state

	defer runningMutex.Unlock()
}

func getRunning() bool {
	return running
}

func addFrameQueue1D(frame *[]Pixel) {
	frameQueue1DMutex.Lock()
	defer frameQueue1DMutex.Unlock()

	frameQueue1D = append(frameQueue1D, frame)
}

func getFrameQueue1D() *[]Pixel {
	frameQueue1DMutex.Lock()
	defer frameQueue1DMutex.Unlock()

	var frame *[]Pixel
	if len(frameQueue1D) > 0 {
		frame = frameQueue1D[0]

		frameQueue1D = frameQueue1D[1:]
	}

	return frame
}

func isEmptyQueue() bool {
	frameQueue1DMutex.Lock()
	defer frameQueue1DMutex.Unlock()

	return len(frameQueue1D) == 0
}

func isWorker() bool {
	global := js.Global()

	workerGlobalScope := global.Get("WorkerGlobalScope")
	if workerGlobalScope.Truthy() {
		return global.InstanceOf(workerGlobalScope)
	}

	return false
}
