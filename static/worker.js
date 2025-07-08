/*
  ===============================================================
  File: load_go.js
  Description: Seperate worker for wasm files
  Author: DryBearr
  ===============================================================
*/

// Predeclare Go object
import "./wasm_exec.js";

const go = new Go();

// Store config
self.computeParams = {};

self.onmessage = (event) => {
  const msg = event.data;
  const type = msg.type;

  if (type === "init") {
    self.computeParams.width = msg.data.width;
    self.computeParams.height = msg.data.height;
    self.computeParams.wasm = msg.data.wasm;

    WebAssembly.instantiateStreaming(
      fetch(self.computeParams.wasm),
      go.importObject,
    )
      .then((result) => {
        go.run(result.instance);
      })
      .catch((error) => {
        self.postMessage({
          type: "log",
          message: `[Worker] error:${error};\n msg: ${JSON.stringify(msg)}\n`,
        });
      });

    self.postMessage({
      type: "log",
      message: `[Worker] width:${msg.data.width}; height:${msg.data.height}, wasm:${msg.data.wasm} initialized`,
    });
  }
};
