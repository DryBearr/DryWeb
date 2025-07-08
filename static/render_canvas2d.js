/*
  ===============================================================
  File: render_canvas2d.js
  Description: Provides utility function for 2d canvas rendering 
  Author: DryBearr
  ===============================================================
*/

const canvas = document.createElement("canvas");

canvas.setAttribute("id", "renderer");

document.querySelector("main").append(canvas);

const canvasCtx = canvas.getContext("2d");

let latestFrame = null;
let width = 0;
let height = 0;
let cancel = false;
let running = false;

export function renderLoop() {
  if (!running) {
    console.log("[RenderLoop] started.");
    running = true;
  }

  if (cancel) {
    running = false;

    console.log("[RenderLoop] canceled.");
    return;
  }

  if (latestFrame) {
    console.log("[RenderLoop] rendering frame.");
    const clamped = new Uint8ClampedArray(latestFrame.buffer);

    const imageData = new ImageData(clamped, width, height);

    if (height !== canvas.height) {
      canvas.height = height;
    }

    if (width !== canvas.width) {
      canvas.width = width;
    }

    canvasCtx.putImageData(imageData, 0, 0);
    console.log("[RenderLoop] rendered correctly i hope.");
  }

  requestAnimationFrame(renderLoop);
}

export function setFrame(frame) {
  if (frame) {
    latestFrame = frame;
  }
}

export function setHeight(h) {
  if (h && h > 0) {
    height = h;
  } else {
    height = 0;
  }
}

export function setWidth(w) {
  if (w && w > 0) {
    width = w;
  } else {
    width = 0;
  }
}

//TODO: better rendering loop manipulation
export function cancelRenderLoop() {
  cancel = true;
}

export function isLoopRunning() {
  return running;
}
