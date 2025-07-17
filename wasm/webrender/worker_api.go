// ===============================================================
// File: worker_api.go
// Description: Provides communication between two workers, implements DryEngine interface
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package webrender

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"
	"time"
	"wasm/dryengine"
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

	size := dryengine.Size{
		Height: params.Get("height").Int(),
		Width:  params.Get("width").Int(),
	}

	Api = &WorkerApi{
		size:      size,
		frameChan: make(chan *dryengine.RenderFrame, 1000), //TODO:
	}

	functionsToRealease := []js.Func{
		js.FuncOf(Api.resizeEventListener),
		js.FuncOf(Api.mouseClickEventListener),
		js.FuncOf(Api.mouseDragEventListener),
		js.FuncOf(Api.mouseDragEndEventListener),
		js.FuncOf(Api.keyDownEventListener),
		js.FuncOf(Api.swipeEventListener),
	}

	for _, f := range functionsToRealease {
		js.Global().Call("addEventListener", "message", f)
		//TODO: prevent memory leak f.Release()
	}
}

type WorkerApi struct {
	resizeHandlers       []dryengine.SizeChangeHandler
	mouseClickHandlers   []dryengine.MouseClickHandler
	mouseDragHandlers    []dryengine.MouseDragHandler
	mouseDragEndHandlers []dryengine.MouseDragEndHandler
	keyDownHandlers      []dryengine.KeyDownHandler
	swipeHandlers        []dryengine.SwipeHandler

	size dryengine.Size

	frameChan chan *dryengine.RenderFrame
}

func (worker *WorkerApi) drawFrame(frame2D *[][]dryengine.Pixel, s dryengine.Size) error {

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

func (worker *WorkerApi) drawFramePartly(frame2D *[][]dryengine.Pixel, s dryengine.Size, c dryengine.Coordinate2D) error {
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

func (worker *WorkerApi) RegisterResizeEventListener(handler dryengine.SizeChangeHandler) error {
	worker.resizeHandlers = append(worker.resizeHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterMouseClickEventListener(handler dryengine.MouseClickHandler) error {
	worker.mouseClickHandlers = append(worker.mouseClickHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterMouseDragEventListener(handler dryengine.MouseDragHandler) error {
	worker.mouseDragHandlers = append(worker.mouseDragHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterMouseDragEndEventListener(handler dryengine.MouseDragEndHandler) error {
	worker.mouseDragEndHandlers = append(worker.mouseDragEndHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterKeyDownEventListener(handler dryengine.KeyDownHandler) error {
	worker.keyDownHandlers = append(worker.keyDownHandlers, handler)
	return nil
}

func (worker *WorkerApi) RegisterSwipeEventListener(handler dryengine.SwipeHandler) error {
	worker.swipeHandlers = append(worker.swipeHandlers, handler)
	return nil
}

func (worker *WorkerApi) GetSize() dryengine.Size {
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

	c := dryengine.Coordinate2D{
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

	c := dryengine.Coordinate2D{
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

	worker.size = dryengine.Size{
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

	c := dryengine.Coordinate2D{
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

	key := dryengine.Key(strings.ToLower(jsKey.String()))
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

func (worker *WorkerApi) swipeEventListener(this js.Value, args []js.Value) any {
	log.Println("[WorkerApi] Recieved swipe event")

	jsObj := worker.getMessageData(args, "swipe")
	if jsObj == nil {
		return nil
	}

	jsSwipeDirection := jsObj.Get("direction")
	if jsSwipeDirection.Type() != js.TypeString {
		log.Println("[WorkerApi] direction is missing or not a string")

		return nil
	}

	swipeDirection := dryengine.Key(strings.ToLower(jsSwipeDirection.String()))
	var direction dryengine.SwipeDirection
	switch swipeDirection {
	case "right":
		direction = dryengine.SwipeRight
	case "left":
		direction = dryengine.SwipeLeft
	case "down":
		direction = dryengine.SwipeDown
	case "up":
		direction = dryengine.SwipeUp
	default:
		log.Println("[WorkerApi] direction is invalid: ", swipeDirection)

		return nil
	}

	for _, handler := range worker.swipeHandlers {
		handler(direction)
	}

	log.Println("[WorkerApi] Notified every swipe handler")

	return nil
}

func (worker *WorkerApi) StartRenderLoop(latency time.Duration) {
	go func() {
		timer := time.NewTimer(latency)
		defer timer.Stop()

		for {
			select {
			case frame, ok := <-worker.frameChan:
				if !ok {
					return
				}

				if frame.C != nil {
					worker.drawFramePartly(frame.Frame, frame.FrameSize, *frame.C)
				} else {
					worker.drawFrame(frame.Frame, frame.FrameSize)
				}

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

func (worker *WorkerApi) AddFrame(renderFrame *dryengine.RenderFrame) {
	worker.frameChan <- renderFrame
}
