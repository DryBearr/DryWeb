/*
  ===============================================================
  File: worker_canvas.js
  Description: Worker canvas for rendering
  Author: DryBearr
  ===============================================================
*/

// TODO: typescript and vite is a pain in one place so for now just gonna use plain javascript :)

self.params = {
  offScreenCanvas: null,
  ctx: null,
};

self.addEventListener("message", (event) => {
  const { type } = event.data;

  switch (type) {
    case "init": {
      console.log("[worker_canvas.js] Received init", event.data);
      self.params.offScreenCanvas = event.data.offScreenCanvas;
      self.params.ctx = self.params.offScreenCanvas.getContext("2d");
      break;
    }

    case "renderFrame": {
      const { pixels, width, height, x, y } = event.data;
      const imageData = new ImageData(
        new Uint8ClampedArray(pixels),
        width,
        height,
      );
      self.params.ctx.putImageData(imageData, x, y);
      break;
    }

    default:
      // Ignore unknown message types
      break;
  }
});
