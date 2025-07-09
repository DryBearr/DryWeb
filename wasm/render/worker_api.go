// ===============================================================
// File: canvas_2d.go
// Description: Provides communicationo between two workers implement Renderer interface
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package render

import (
	"log"
	"syscall/js"
)

var (
	Api *WorkerApi
)

func init() {
	workerGlobalScope := js.Global().Get("WorkerGlobalScope")
	if workerGlobalScope.Truthy() && !js.Global().InstanceOf(workerGlobalScope) {
		panic("can't run outside worker!!!")
	}

	params := js.Global().Get("computeParams")
	if !params.Truthy() {
		panic("computeParams is not defined in the worker")
	}

	size := Size{
		Height: params.Get("height").Int(),
		Width:  params.Get("width").Int(),
	}

	Api = &WorkerApi{
		size: size,
	}

	js.Global().Call("addEventListener", "message", js.FuncOf(Api.resizeEventListener))
}

type WorkerApi struct {
	handlers []SizeChangeHandler

	size Size

	//TODO:
	//pixelBuf      []byte
	//pixelBufMutex sync.Mutex
}

func (worker *WorkerApi) DrawFrame(frame *[]Pixel, size Size) {

	//TODO: check sizes
	pixelBuf := make([]byte, len(*frame)*4)

	for idx, pixel := range *frame {
		i := idx * 4
		pixelBuf[i+0] = pixel.R
		pixelBuf[i+1] = pixel.G
		pixelBuf[i+2] = pixel.B
		pixelBuf[i+3] = pixel.A
	}

	// Create a JS Uint8Array from your []byte
	uint8Array := js.Global().Get("Uint8Array").New(len(pixelBuf))
	js.CopyBytesToJS(uint8Array, pixelBuf)

	// Create a JS object to hold message
	msg := js.Global().Get("Object").New()
	msg.Set("type", "pixels")
	msg.Set("pixels", uint8Array)
	msg.Set("width", size.Width)
	msg.Set("height", size.Height)

	// Send
	js.Global().Call("postMessage", msg)
}

func (worker *WorkerApi) RegisterResizeEventListener(handler SizeChangeHandler) {
	worker.handlers = append(worker.handlers, handler)
}

func (worker *WorkerApi) GetSize() Size {
	return worker.size
}

func (worker *WorkerApi) resizeEventListener(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return nil
	}

	jsObj := args[0].Get("data")

	if jsObj.Type() != js.TypeObject {
		return nil
	}

	messageType := jsObj.Get("type")
	if messageType.Type() != js.TypeString {
		return nil
	}

	if messageType.String() != "resize" {
		return nil
	}

	log.Println("[WorkerApi] Recieved resize event")

	widthVal := jsObj.Get("width")
	if widthVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] width is missing or not a number")

		return nil
	}

	heightVal := jsObj.Get("height")
	if heightVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] height is missing or not a number")

		return nil
	}

	worker.size = Size{
		Width:  widthVal.Int(),
		Height: heightVal.Int(),
	}

	for _, handler := range worker.handlers {
		handler(worker.size)
	}

	log.Println("[WorkerApi] Notified every handler")

	return nil
}
