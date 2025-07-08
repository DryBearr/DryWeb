/*
  ===============================================================
  File: core.js
  Description: Init worker base on given url params
  Author: DryBearr
  ===============================================================
*/

import {
  isLoopRunning,
  renderLoop,
  setFrame,
  setWidth,
  setHeight,
} from "./render_canvas2d.js";
import config from "./available_wasms.json" with { type: "json" };

//Config
let loadWasm = config.default_wasm_path;

//Try to load wasm base on given param in url
const queryString = new URLSearchParams(window.location.search);

if (queryString.has("wasm")) {
  const wasmParam = queryString.get("wasm");

  if (wasmParam in config.available_wasms) {
    loadWasm = config.available_wasms[wasmParam];
  }
}

console.log("[Loader] Using wasm:", loadWasm);

//Set start width and height
const width =
  document.querySelector("main")?.offsetWidth ||
  document.querySelector("body")?.offsetWidth;
const height =
  document.querySelector("main")?.offsetHeight ||
  document.querySelector("body")?.offsetHeight;

if (!width || !height) {
  console.error(
    "[Core] no body or main element, can't get height & width. bro wtf ...",
  );
}

console.log("[Core] Init worker...");
// Init worker for wasm
const worker = new Worker("worker.js", { type: "module" });

worker.postMessage({
  type: "init",
  data: {
    wasm: loadWasm,
    width: width,
    height: height,
  },
});

console.log("[Core] Init worker Done.");

worker.addEventListener("message", function (event) {
  const data = event.data;
  if (data.type === "log") {
    console.log(data.message);
  }
});

worker.addEventListener("message", function (event) {
  console.log("[Core] reciving pixels from worker.");
  const data = event.data;
  if (data.type === "pixels") {
    if (!data.height || !data.width || !data.pixels) {
      return;
    }

    if (!isLoopRunning()) {
      console.log("[Core] started render loop.");
      renderLoop();
    }

    setWidth(data.width);
    setHeight(data.height);
    setFrame(data.pixels);
  }
});
