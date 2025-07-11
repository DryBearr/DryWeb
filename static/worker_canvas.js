/*
  ===============================================================
  File: worker_api.js
  Description: Worker canvas rendering 
  Author: DryBearr
  ===============================================================
*/

let loopId = null;

self.params = {
  offScreenCanvas: null,
  ctx: null,
};

self.addEventListener("message", (event) => {
  const type = event.data.type;

  if (type === "init") {
    console.log("[worker_canvas.js] Received init", event.data);
    self.params.offScreenCanvas = event.data.offScreenCanvas;
    self.params.ctx = self.params.offScreenCanvas.getContext("2d");
  }
});

self.addEventListener("message", async (event) => {
  if (event.data.type === "frame") {
    const { pixels, width, height } = event.data;

    const imageData = new ImageData(
      new Uint8ClampedArray(pixels),
      width,
      height,
    );

    if (
      self.params.offScreenCanvas.width !== width ||
      self.params.offScreenCanvas.height !== height
    ) {
      self.params.offScreenCanvas.width = width;
      self.params.offScreenCanvas.height = height;
    }

    self.params.ctx.putImageData(imageData, 0, 0);
  }
});

self.addEventListener("message", async (event) => {
  if (event.data.type === "framePart") {
    const { pixels, width, height, x, y } = event.data;

    const imageData = new ImageData(
      new Uint8ClampedArray(pixels),
      width,
      height,
    );

    self.params.ctx.putImageData(imageData, x, y);
  }
});
