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
    self.params.ctx = self.params.offScreenCanvas.getContext("bitmaprenderer");
  }
});

self.addEventListener("message", async (event) => {
  if (event.data.type === "pixels") {
    const { pixels, width, height } = event.data;

    const imageData = new ImageData(
      new Uint8ClampedArray(pixels),
      width,
      height,
    );

    const bitmap = await createImageBitmap(imageData);

    if (
      self.params.offScreenCanvas.width !== width ||
      self.params.offScreenCanvas.height !== height
    ) {
      self.params.offScreenCanvas.width = width;
      self.params.offScreenCanvas.height = height;
    }

    self.params.ctx.transferFromImageBitmap(bitmap);

    bitmap.close();
  }
});
