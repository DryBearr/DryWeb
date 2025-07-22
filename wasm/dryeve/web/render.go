// ===============================================================
// File: render.go
// Description: Implements render.Renderer interface for web
// Author: DryBearr
// ===============================================================

//go:build js && wasm

package web

import (
	"fmt"
	"syscall/js"
	"wasm/dryeve/models"
	"wasm/dryeve/render"
	"wasm/dryeve/util"
)

type WebRenderer struct{}

func NewWebRenderer() render.Renderer {
	return &WebRenderer{}
}

func (r *WebRenderer) RenderRect(rect models.Rect, pixel models.Pixel) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("RenderRect failed: %v", rec)
		}
	}()

	msg := js.Global().Get("Object").New()
	msg.Set("type", "renderRect")
	msg.Set("x", rect.C.X)
	msg.Set("y", rect.C.Y)
	msg.Set("width", rect.Width)
	msg.Set("height", rect.Height)
	msg.Set("color", util.EncodeColorHex(pixel))
	js.Global().Call("postMessage", msg)

	return nil
}

func (r *WebRenderer) RenderCircle(circle models.Circle, pixel models.Pixel) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("RenderCircle failed: %v", rec)
		}
	}()

	msg := js.Global().Get("Object").New()
	msg.Set("type", "renderCircle")
	msg.Set("x", circle.Center.X)
	msg.Set("y", circle.Center.Y)
	msg.Set("radius", circle.R)
	msg.Set("startAngle", circle.StartAngle)
	msg.Set("endAngle", circle.EndAngle)
	msg.Set("color", util.EncodeColorHex(pixel))
	js.Global().Call("postMessage", msg)

	return nil
}

func (r *WebRenderer) RenderLine(line models.Line, pixel models.Pixel) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("RenderLine failed: %v", rec)
		}
	}()

	msg := js.Global().Get("Object").New()
	msg.Set("type", "renderLine")
	msg.Set("startX", line.Start.X)
	msg.Set("startY", line.Start.Y)
	msg.Set("endX", line.End.X)
	msg.Set("endY", line.End.Y)
	msg.Set("width", line.Width)
	msg.Set("color", util.EncodeColorHex(pixel))
	js.Global().Call("postMessage", msg)

	return nil
}

func (r *WebRenderer) RenderPixel(point models.Point2D, pixel models.Pixel) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("RenderPixel failed: %v", rec)
		}
	}()

	msg := js.Global().Get("Object").New()
	msg.Set("type", "renderPixel")
	msg.Set("x", point.X)
	msg.Set("y", point.Y)
	msg.Set("color", util.EncodeColorHex(pixel))
	js.Global().Call("postMessage", msg)

	return nil
}

func (r *WebRenderer) RenderFrame(renderFrame models.RenderFrame) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("RenderFrame failed: %v", rec)
		}
	}()

	if renderFrame.Frame == nil {
		return fmt.Errorf("RenderFrame failed: Frame is nil")
	}

	if len(*renderFrame.Frame) == 0 {
		return nil
	}

	width := len((*renderFrame.Frame)[0])
	height := len((*renderFrame.Frame))

	pixelBuf := make([]byte, width*height*4)

	for y, frame1D := range *renderFrame.Frame {
		for x, pixel := range frame1D {
			i := (y*width + x) * 4
			pixelBuf[i+0] = pixel.R
			pixelBuf[i+1] = pixel.G
			pixelBuf[i+2] = pixel.B
			pixelBuf[i+3] = pixel.A
		}
	}

	uint8Array := js.Global().Get("Uint8Array").New(len(pixelBuf))
	js.CopyBytesToJS(uint8Array, pixelBuf)

	x := 0
	y := 0

	if renderFrame.C != nil {
		x = int(renderFrame.C.X)
		y = int(renderFrame.C.Y)
	}

	msg := js.Global().Get("Object").New()
	msg.Set("type", "renderFrame")
	msg.Set("pixels", uint8Array)
	msg.Set("width", width)
	msg.Set("height", height)
	msg.Set("x", x)
	msg.Set("y", y)
	js.Global().Call("postMessage", msg)

	return nil
}
