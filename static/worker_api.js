/*
  ===============================================================
  File: worker_api.js
  Description: Worker for wasm worker_api.go and wasm that relies on the api 
  Author: DryBearr
  ===============================================================
*/

// Predeclare Go object
import "./wasm_exec.js";

const go = new Go();

// Store config
self.computeParams = {};

self.addEventListener("message", (event) => {
  const type = event.data.type;
  if (type === "init") {
    const { width, height, wasm } = event.data;

    self.computeParams.width = width;
    self.computeParams.height = height;
    self.computeParams.wasm = wasm;

    WebAssembly.instantiateStreaming(
      fetch(self.computeParams.wasm),
      go.importObject,
    )
      .then((result) => {
        go.run(result.instance);
      })
      .catch((error) => {
        console.error(
          `[WorkerApiJS] error:${error};\n msg: ${JSON.stringify(event.data)}\n`,
        );
      });

    console.log(
      `[WorkerApiJS] Received init message. width: ${width}, height: ${height}, wasm: ${wasm}`,
    );
  }
});
