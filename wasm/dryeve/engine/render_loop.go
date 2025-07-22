// ===============================================================
// File: render_loop.go
// Description: Runs the fixed-interval render loop for DryEve engine.
// Author: DryBearr
// ===============================================================

package engine

import (
	"time"
	"wasm/dryeve/events"
	"wasm/dryeve/models"
	"wasm/dryeve/render"
)

type Engine struct {
	Renderer render.Renderer
	Events   events.Events

	latency time.Duration

	frameChan chan models.RenderFrame
}

func NewEngine(renderer render.Renderer, events events.Events, latency time.Duration, frameBuffSize int) *Engine {
	return &Engine{
		Renderer:  renderer,
		Events:    events,
		latency:   latency,
		frameChan: make(chan models.RenderFrame, frameBuffSize),
	}
}

// TODO: change rendering logic to support other renderer features
func (engine *Engine) StartRenderLoop() {
	go func() {
		timer := time.NewTimer(engine.latency)
		defer timer.Stop()

		for {
			select {
			case frame, ok := <-engine.frameChan:
				if !ok {
					return
				}

				engine.Renderer.RenderFrame(frame)

				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(engine.latency)

			case <-timer.C:
				timer.Reset(engine.latency)
			}
		}
	}()
}

func (engine *Engine) AddFrame(renderFrame models.RenderFrame) {
	go func() {
		engine.frameChan <- renderFrame
	}()
}
