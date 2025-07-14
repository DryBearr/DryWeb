// ===============================================================
// File: rendering.go
// Description: Logic of rendering for the game of life
// Author: DryBearr
// ===============================================================

package core

import (
	"time"
)

func StartRenderLoop() {
	go func() {
		timer := time.NewTimer(latency)
		defer timer.Stop()

		for {
			select {
			case frame, ok := <-frameChan:
				if !ok {
					return
				}

				if frame.C != nil {
					api.DrawFramePartly(frame.Frame, frame.Size, *frame.C)
				} else {
					api.DrawFrame(frame.Frame, frame.Size)
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

func AddFrame(f RenderFrame) {
	frameChan <- f
}
