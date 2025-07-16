// ===============================================================
// File: worker_api.go
// Description: Provides communicationo between two workers, implements Renderer interface
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package webrender

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"
	"wasm/render"
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

	size := render.Size{
		Height: params.Get("height").Int(),
		Width:  params.Get("width").Int(),
	}

	Api = &WorkerApi{
		size: size,
	}

	js.Global().Call("addEventListener", "message", js.FuncOf(Api.resizeEventListener))
	js.Global().Call("addEventListener", "message", js.FuncOf(Api.mouseClickEventListener))
	js.Global().Call("addEventListener", "message", js.FuncOf(Api.mouseDragEventListener))
	js.Global().Call("addEventListener", "message", js.FuncOf(Api.mouseDragEndEventListener))
	js.Global().Call("addEventListener", "message", js.FuncOf(Api.keyDownEventListener))
}

type WorkerApi struct {
	resizeHandlers       []render.SizeChangeHandler
	mouseClickHandlers   []render.MouseClickHandler
	mouseDragHandlers    []render.MouseDragHandler
	mouseDragEndHandlers []render.MouseDragEndHandler
	keyDownHandlers      []render.KeyDownHandler

	size render.Size

	//TODO:
	//pixelBuf      []byte
	//pixelBufMutex sync.Mutex
}

func (worker *WorkerApi) DrawFrame(frame2D *[][]render.Pixel, s render.Size) error {

	//TODO: check sizes
	pixelBuf := make([]byte, s.Width*s.Height*4)

	for y, frame1D := range *frame2D {
		for x, pixel := range frame1D {
			i := (y*s.Width + x) * 4
			pixelBuf[i+0] = pixel.R
			pixelBuf[i+1] = pixel.G
			pixelBuf[i+2] = pixel.B
			pixelBuf[i+3] = pixel.A
		}
	}

	uint8Array := js.Global().Get("Uint8Array").New(len(pixelBuf))
	js.CopyBytesToJS(uint8Array, pixelBuf)

	msg := js.Global().Get("Object").New()
	msg.Set("type", "frame")
	msg.Set("pixels", uint8Array)
	msg.Set("width", s.Width)
	msg.Set("height", s.Height)

	js.Global().Call("postMessage", msg)

	return nil
}

func (worker *WorkerApi) DrawFramePartly(frame2D *[][]render.Pixel, s render.Size, c render.Coordinate) error {
	if !worker.size.EqualOrGreater(s) {
		return fmt.Errorf("invalid size")
	}

	//TODO: check sizes
	pixelBuf := make([]byte, s.Width*s.Height*4)

	for y, frame1D := range *frame2D {
		for x, pixel := range frame1D {
			i := (y*s.Width + x) * 4
			pixelBuf[i+0] = pixel.R
			pixelBuf[i+1] = pixel.G
			pixelBuf[i+2] = pixel.B
			pixelBuf[i+3] = pixel.A
		}
	}

	uint8Array := js.Global().Get("Uint8Array").New(len(pixelBuf))
	js.CopyBytesToJS(uint8Array, pixelBuf)

	msg := js.Global().Get("Object").New()
	msg.Set("type", "framePart")
	msg.Set("pixels", uint8Array)
	msg.Set("width", s.Width)
	msg.Set("height", s.Height)
	msg.Set("x", c.X)
	msg.Set("y", c.Y)

	js.Global().Call("postMessage", msg)

	return nil
}

func (worker *WorkerApi) RegisterResizeEventListener(handler render.SizeChangeHandler) error {
	worker.resizeHandlers = append(worker.resizeHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterMouseClickEventListener(handler render.MouseClickHandler) error {
	worker.mouseClickHandlers = append(worker.mouseClickHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterMouseDragEventListener(handler render.MouseDragHandler) error {
	worker.mouseDragHandlers = append(worker.mouseDragHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterMouseDragEndEventListener(handler render.MouseDragEndHandler) error {
	worker.mouseDragEndHandlers = append(worker.mouseDragEndHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterKeyDownEventListener(handler render.KeyDownHandler) error {
	worker.keyDownHandlers = append(worker.keyDownHandlers, handler)
	return nil
}

func (worker *WorkerApi) GetSize() render.Size {
	return worker.size
}

func (worker *WorkerApi) getMessageData(args []js.Value, t string) *js.Value {
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

	if messageType.String() != t {
		return nil
	}

	return &jsObj
}

func (worker *WorkerApi) mouseDragEventListener(this js.Value, args []js.Value) any {
	log.Println("[WorkerApi] Recieved mouseDrag event")

	jsObj := worker.getMessageData(args, "mouseDrag")
	if jsObj == nil {
		return nil
	}

	xVal := jsObj.Get("x")
	if xVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] x is missing or not a number")

		return nil
	}

	yVal := jsObj.Get("y")
	if yVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] y is missing or not a number")

		return nil
	}

	c := render.Coordinate{
		X: xVal.Int(),
		Y: yVal.Int(),
	}

	for _, handler := range worker.mouseDragHandlers {
		handler(c)
	}

	log.Println("[WorkerApi] Notified every mouseDrag handler")

	return nil
}

func (worker *WorkerApi) mouseClickEventListener(this js.Value, args []js.Value) any {
	log.Println("[WorkerApi] Recieved mouseClick event")

	jsObj := worker.getMessageData(args, "mouseClick")
	if jsObj == nil {
		return nil
	}

	xVal := jsObj.Get("x")
	if xVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] x is missing or not a number")

		return nil
	}

	yVal := jsObj.Get("y")
	if yVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] y is missing or not a number")

		return nil
	}

	c := render.Coordinate{
		X: xVal.Int(),
		Y: yVal.Int(),
	}

	for _, handler := range worker.mouseClickHandlers {
		handler(c)
	}

	log.Println("[WorkerApi] Notified every mouseClick handler")

	return nil

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

	worker.size = render.Size{
		Width:  widthVal.Int(),
		Height: heightVal.Int(),
	}

	for _, handler := range worker.resizeHandlers {
		handler(worker.size)
	}

	log.Println("[WorkerApi] Notified every resize handler")

	return nil
}

func (worker *WorkerApi) mouseDragEndEventListener(this js.Value, args []js.Value) any {
	log.Println("[WorkerApi] Recieved mouseDragEnd event")

	jsObj := worker.getMessageData(args, "mouseDragEnd")
	if jsObj == nil {
		return nil
	}

	xVal := jsObj.Get("x")
	if xVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] x is missing or not a number")

		return nil
	}

	yVal := jsObj.Get("y")
	if yVal.Type() != js.TypeNumber {
		log.Println("[WorkerApi] y is missing or not a number")

		return nil
	}

	c := render.Coordinate{
		X: xVal.Int(),
		Y: yVal.Int(),
	}

	for _, handler := range worker.mouseDragEndHandlers {
		handler(c)
	}

	log.Println("[WorkerApi] Notified every mouseDragEnd handler")

	return nil

}

// TODO: key down, key up
func (worker *WorkerApi) keyDownEventListener(this js.Value, args []js.Value) any {
	log.Println("[WorkerApi] Recieved keyDown event")

	jsObj := worker.getMessageData(args, "keyDown")
	if jsObj == nil {
		return nil
	}

	jsKey := jsObj.Get("key")
	if jsKey.Type() != js.TypeString {
		log.Println("[WorkerApi] key is missing or not a string")

		return nil
	}

	key := render.Key(strings.ToLower(jsKey.String()))
	if len(key) > 1 {
		log.Println("[WorkerApi] key is invalid: ", jsKey.String())

		return nil
	}

	for _, handler := range worker.keyDownHandlers {
		handler(key)
	}

	log.Println("[WorkerApi] Notified every keyDown handler")

	return nil
}
