/*
  ===============================================================
  File: worker_api.js
  Description: DryEngine web relies on this worker 
  Author: DryBearr
  ===============================================================
*/

// TODO: typescript and vite is a pain in one place so for now just gonna use plain javascript :)

import "./wasm_exec.js";

const go = new Go();

self.computeParams = {};

self.addEventListener("message", (event) => {
  const { type } = event.data;

  switch (type) {
    case "init": {
      const { width, height, wasm } = event.data;

      self.computeParams = { width, height, wasm };

      console.log(
        `[WorkerApiJS] Received init message. width: ${width}, height: ${height}, wasm: ${wasm}`,
      );

      WebAssembly.instantiateStreaming(fetch(wasm), go.importObject)
        .then((result) => {
          go.run(result.instance);
        })
        .catch((error) => {
          console.error(
            `[WorkerApiJS] error: ${error};\n msg: ${JSON.stringify(event.data)}\n`,
          );
        });

      break;
    }

    default:
      // No other message types handled yet
      break;
  }
});
