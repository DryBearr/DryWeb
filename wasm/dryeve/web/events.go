// ===============================================================
// File: events.go
// Description: Implements Events interface
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package web

import (
	"strings"
	"syscall/js"
	"wasm/dryeve/events"
	"wasm/dryeve/models"
)

type WebEvents struct {
	resizeHandlers       []models.SizeChangeHandler
	mouseClickHandlers   []models.MouseClickHandler
	mouseDragHandlers    []models.MouseDragHandler
	mouseDragEndHandlers []models.MouseDragEndHandler
	keyDownHandlers      []models.KeyDownHandler
	swipeHandlers        []models.SwipeHandler
}

func NewWebEvents() events.Events {
	webEvents := &WebEvents{}

	functionsToRealease := []js.Func{
		js.FuncOf(webEvents.resizeEventListener),
		js.FuncOf(webEvents.mouseClickEventListener),
		js.FuncOf(webEvents.mouseDragEventListener),
		js.FuncOf(webEvents.mouseDragEndEventListener),
		js.FuncOf(webEvents.keyDownEventListener),
		js.FuncOf(webEvents.swipeEventListener),
	}

	for _, f := range functionsToRealease {
		js.Global().Call("addEventListener", "message", f)
		//TODO: prevent memory leak f.Release()
	}

	return webEvents
}

func (e *WebEvents) RegisterResizeEventListener(handler models.SizeChangeHandler) error {
	e.resizeHandlers = append(e.resizeHandlers, handler)

	return nil
}

func (e *WebEvents) RegisterMouseClickEventListener(handler models.MouseClickHandler) error {
	e.mouseClickHandlers = append(e.mouseClickHandlers, handler)

	return nil
}

func (e *WebEvents) RegisterMouseDragEventListener(handler models.MouseDragHandler) error {
	e.mouseDragHandlers = append(e.mouseDragHandlers, handler)

	return nil
}

func (e *WebEvents) RegisterMouseDragEndEventListener(handler models.MouseDragEndHandler) error {
	e.mouseDragEndHandlers = append(e.mouseDragEndHandlers, handler)
	return nil
}

func (e *WebEvents) RegisterKeyDownEventListener(handler models.KeyDownHandler) error {
	e.keyDownHandlers = append(e.keyDownHandlers, handler)

	return nil
}

func (e *WebEvents) RegisterSwipeEventListener(handler models.SwipeHandler) error {
	e.swipeHandlers = append(e.swipeHandlers, handler)

	return nil
}

func (e *WebEvents) getMessageData(args []js.Value, t string) *js.Value {
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

func (e *WebEvents) mouseDragEventListener(this js.Value, args []js.Value) any {
	//TODO: recover from error and some how log it

	jsObj := e.getMessageData(args, "mouseDrag")
	if jsObj == nil {
		return nil
	}

	xVal := jsObj.Get("x")
	if xVal.Type() != js.TypeNumber {
		return nil
	}

	yVal := jsObj.Get("y")
	if yVal.Type() != js.TypeNumber {
		return nil
	}

	c := models.Point2D{
		X: float32(xVal.Int()),
		Y: float32(yVal.Int()),
	}

	for _, handler := range e.mouseDragHandlers {
		handler(c)
	}

	return nil
}

func (e *WebEvents) mouseClickEventListener(this js.Value, args []js.Value) any {
	//TODO: recover from error and some how log it

	jsObj := e.getMessageData(args, "mouseClick")
	if jsObj == nil {
		return nil
	}

	xVal := jsObj.Get("x")
	if xVal.Type() != js.TypeNumber {
		return nil
	}

	yVal := jsObj.Get("y")
	if yVal.Type() != js.TypeNumber {
		return nil
	}

	c := models.Point2D{
		X: float32(xVal.Int()),
		Y: float32(yVal.Int()),
	}

	for _, handler := range e.mouseClickHandlers {
		handler(c)
	}

	return nil
}

func (e *WebEvents) resizeEventListener(this js.Value, args []js.Value) any {
	//TODO: recover from error and some how log it

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

	widthVal := jsObj.Get("width")
	if widthVal.Type() != js.TypeNumber {
		return nil
	}

	heightVal := jsObj.Get("height")
	if heightVal.Type() != js.TypeNumber {
		return nil
	}

	for _, handler := range e.resizeHandlers {
		handler(widthVal.Int(), heightVal.Int())
	}

	return nil
}

func (e *WebEvents) mouseDragEndEventListener(this js.Value, args []js.Value) any {
	//TODO: recover from error and some how log it

	jsObj := e.getMessageData(args, "mouseDragEnd")
	if jsObj == nil {
		return nil
	}

	xVal := jsObj.Get("x")
	if xVal.Type() != js.TypeNumber {
		return nil
	}

	yVal := jsObj.Get("y")
	if yVal.Type() != js.TypeNumber {
		return nil
	}

	c := models.Point2D{
		X: float32(xVal.Int()),
		Y: float32(yVal.Int()),
	}

	for _, handler := range e.mouseDragEndHandlers {
		handler(c)
	}

	return nil
}

// TODO: key down, key up
func (e *WebEvents) keyDownEventListener(this js.Value, args []js.Value) any {
	//TODO: recover from error and some how log it

	jsObj := e.getMessageData(args, "keyDown")
	if jsObj == nil {
		return nil
	}

	jsKey := jsObj.Get("key")
	if jsKey.Type() != js.TypeString {
		return nil
	}

	s := strings.ToLower(jsKey.String())

	key := models.KeyUnknown
	switch strings.ToLower(s) {
	case "a":
		key = models.KeyA
	case "w":
		key = models.KeyW
	case "s":
		key = models.KeyS
	case "d":
		key = models.KeyD
	}

	for _, handler := range e.keyDownHandlers {
		handler(key)
	}

	return nil
}

func (e *WebEvents) swipeEventListener(this js.Value, args []js.Value) any {
	//TODO: recover from error and some how log it

	jsObj := e.getMessageData(args, "swipe")
	if jsObj == nil {
		return nil
	}

	jsSwipeDirection := jsObj.Get("direction")
	if jsSwipeDirection.Type() != js.TypeString {
		return nil
	}

	swipeDirection := strings.ToLower(jsSwipeDirection.String())
	var direction models.SwipeDirection
	switch swipeDirection {
	case "right":
		direction = models.SwipeRight
	case "left":
		direction = models.SwipeLeft
	case "down":
		direction = models.SwipeDown
	case "up":
		direction = models.SwipeUp
	default:
		return nil
	}

	for _, handler := range e.swipeHandlers {
		handler(direction)
	}

	return nil
}
