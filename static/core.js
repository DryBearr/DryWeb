/*
  ===============================================================
  File: core.js
  Description: Core script for website inserting links, initilizing workers, establishint communication between them 
  Author: DryBearr
  ===============================================================
*/

//Config
import config from "./available_wasms.json" with { type: "json" };
let loadWasm = config.default_wasm_path;

//Insert links
const nav = document.getElementById("wasm-links");

const urlParams = new URLSearchParams(window.location.search);
const currentHref = urlParams.get("wasm") || config.default_wasm;

Object.entries(config.available_wasms).forEach(([key, path]) => {
  const link = document.createElement("a");
  link.innerText = key;
  link.href = `?wasm=${key}`;
  if (key === currentHref) {
    link.classList.add("active-link");
  }
  nav.append(link);
});

//Try to load wasm base on given param in url
const queryString = new URLSearchParams(window.location.search);

if (queryString.has("wasm")) {
  const wasmParam = queryString.get("wasm");

  if (wasmParam in config.available_wasms) {
    loadWasm = config.available_wasms[wasmParam];
  }
}

console.log("[Core] Using wasm:", loadWasm);

// Init Canvas and off screen canvas
const canvas = document.createElement("canvas");
canvas.setAttribute("id", "renderer");

const canvasParent =
  document.querySelector("main") || document.querySelector("body");
if (!canvasParent) {
  console.error("[Core] no body or main element. bro wtf ...");
}

const width = canvasParent.offsetWidth;
const height = canvasParent.offsetHeight;

if (!width || !height) {
  console.error("[Core] can't get height & width. bro wtf ...");
}

canvas.setAttribute("width", width);
canvas.setAttribute("height", height);
canvasParent.append(canvas);

const offScreenCanvas = canvas.transferControlToOffscreen();

// Init worker for wasm
console.log("[Core] Init workers...");

const workerApi = new Worker("worker_api.js", { type: "module" });
const workerCanvas = new Worker("worker_canvas.js", { type: "module" });

workerApi.postMessage({
  type: "init",
  wasm: loadWasm,
  width: width,
  height: height,
});

workerCanvas.postMessage(
  {
    type: "init",
    offScreenCanvas: offScreenCanvas,
  },
  [offScreenCanvas],
);

workerApi.addEventListener("message", function (event) {
  const data = event.data;
  if (data.type === "pixels") {
    workerCanvas.postMessage({
      ...data,
    });
  }
});

const resizeObserver = new ResizeObserver((entries) => {
  for (let entry of entries) {
    const width = entry.contentRect.width;
    const height = entry.contentRect.height;

    console.log("[Core] Resized to", width, height);

    workerApi.postMessage({
      type: "resize",
      width: Math.floor(width),
      height: Math.floor(height),
    });
  }
});

resizeObserver.observe(canvasParent);

console.log("[Core] Init workers Done.");
