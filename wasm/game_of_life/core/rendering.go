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
		timer := time.NewTimer(Latency)
		defer timer.Stop()

		for {
			select {
			case frame, ok := <-FrameChan:
				if !ok {
					return
				}

				if frame.C != nil {
					Api.DrawFramePartly(frame.Frame, frame.Size, *frame.C)
				} else {
					Api.DrawFrame(frame.Frame, frame.Size)
				}

				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(Latency)

			case <-timer.C:
				timer.Reset(Latency)
			}
		}
	}()
}

func AddFrame(f RenderFrame) {
	FrameChan <- f
}
