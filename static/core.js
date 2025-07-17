/*
  ===============================================================
  File: core.js
  Description: Core script for website inserting links, initilizing workers, establishint communication between them 
  Author: DryBearr
  ===============================================================
*/

//Config
import config from "./config.json" with { type: "json" };

//Insert links
const nav = document.getElementById("wasm-links");

const urlParams = new URLSearchParams(window.location.search);
const currentHref = urlParams.get("wasm") || config.default_wasm;

Object.entries(config.available_wasms).forEach(([key, path]) => {
  const link = document.createElement("a");
  const span = document.createElement("span");
  span.textContent = key;
  link.href = `?wasm=${key}`;
  link.appendChild(span);
  if (key === currentHref) {
    link.classList.add("active-link");
  }
  nav.append(link);
});

//Try to load wasm base on given param in url
const queryString = new URLSearchParams(window.location.search);
let loadWasm = config.default_wasm_path;

if (queryString.has("wasm")) {
  const wasmParam = queryString.get("wasm");

  if (wasmParam in config.available_wasms) {
    loadWasm = config.available_wasms[wasmParam];
  }
}

console.log("[Core] Using wasm:", loadWasm);

//Create Window for renderer

let { width, height } = config.window_size.default;
if (isMobileViewport()) {
  ({ width, height } = config.window_size.mobile);
}

const anchor = document.querySelector("main") || document.querySelector("body");
if (!anchor) {
  console.error("[Core] can't find main or body element wtf bro.");
}

const windowDiv = document.createElement("div");
windowDiv.setAttribute("class", "render-window");
anchor.append(windowDiv);

// Init Canvas and off screen canvas
const canvas = document.createElement("canvas");
canvas.setAttribute("id", "renderer");

canvas.setAttribute("width", width);
canvas.setAttribute("height", height);
windowDiv.append(canvas);

const offScreenCanvas = canvas.transferControlToOffscreen();

// Init worker for wasm
console.log("[Core] Init workers...");

let workerApi = new Worker("worker_api.js", { type: "module" });
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
  if (data.type === "frame") {
    workerCanvas.postMessage({
      ...data,
    });
  }
});

workerApi.addEventListener("message", function (event) {
  const data = event.data;
  if (data.type === "framePart") {
    workerCanvas.postMessage({
      ...data,
    });
  }
});

//Controls

const controlsDiv = document.createElement("div");
controlsDiv.setAttribute("class", "controls");

anchor.append(controlsDiv);

//Reload Wasm
const reloadWasmButton = document.createElement("button");
reloadWasmButton.textContent = "Reload";
reloadWasmButton.addEventListener("click", () => {
  workerApi.terminate();

  workerApi = new Worker("./worker_api.js", { type: "module" });

  workerApi.postMessage({
    type: "init",
    wasm: loadWasm,
    width: width,
    height: height,
  });

  workerApi.addEventListener("message", function (event) {
    const data = event.data;
    if (data.type === "frame") {
      workerCanvas.postMessage({
        ...data,
      });
    }
  });

  workerApi.addEventListener("message", function (event) {
    const data = event.data;
    if (data.type === "framePart") {
      workerCanvas.postMessage({
        ...data,
      });
    }
  });
});
reloadWasmButton.setAttribute("class", "reload-button");

controlsDiv.append(reloadWasmButton);

//On Canvas Drag event logic
const getCanvasCoordinates = (event) => {
  const rect = canvas.getBoundingClientRect();
  if (event.touches && event.touches.length > 0) {
    return {
      x: Math.floor(event.touches[0].clientX - rect.left),
      y: Math.floor(event.touches[0].clientY - rect.top),
    };
  } else {
    return {
      x: Math.floor(event.clientX - rect.left),
      y: Math.floor(event.clientY - rect.top),
    };
  }
};

let isDraging = false;
let prevPoint = null;

const handleDragStart = (event) => {
  event.preventDefault();

  isDraging = true;

  const { x, y } = getCanvasCoordinates(event);

  prevPoint = { x, y };
};

const handleDragMove = (event) => {
  event.preventDefault();
  if (!isDraging) return;

  const { x, y } = getCanvasCoordinates(event);

  if (prevPoint != null) {
    workerApi.postMessage({
      type: "mouseDrag",
      x: prevPoint.x,
      y: prevPoint.y,
    });

    prevPoint = null;
  }

  workerApi.postMessage({ type: "mouseDrag", x, y });
};

const handleDragEnd = (event) => {
  event.preventDefault();

  const { x, y } = getCanvasCoordinates(event);

  if (prevPoint) {
    const dx = Math.abs(x - prevPoint.x);
    const dy = Math.abs(y - prevPoint.y);
    if (dx < 3 && dy < 3) {
      workerApi.postMessage({
        type: "mouseClick",
        x: prevPoint.x,
        y: prevPoint.y,
      });
    }
  } else {
    workerApi.postMessage({
      type: "mouseDragEnd",
      x,
      y,
    });
  }

  isDraging = false;
  prevPoint = null;
};

canvas.addEventListener("mousedown", (event) => handleDragStart(event));
canvas.addEventListener("mousemove", (event) => handleDragMove(event));
canvas.addEventListener("mouseup", (event) => handleDragEnd(event));

canvas.addEventListener("touchstart", (event) => handleDragStart(event));
canvas.addEventListener("touchmove", (event) => handleDragMove(event));
canvas.addEventListener("touchend", (event) => handleDragEnd(event));

console.log("[Core] Init workers Done.");

//Event key
//TODO: window or document hmmm
// window.addEventListener("keydown", (event) => {
//   workerApi.postMessage({
//     type: "keyDown",
//     key: event.key,
//   });
// });
document.addEventListener("keydown", (event) => {
  workerApi.postMessage({
    type: "keyDown",
    key: event.key,
  });
});

//Swipe Event

let touchStartX = 0;
let touchStartY = 0;
let touchEndX = 0;
let touchEndY = 0;

const swipeZone = document.getElementById("renderer");

swipeZone.addEventListener("touchstart", function (e) {
  touchStartX = e.changedTouches[0].screenX;
  touchStartY = e.changedTouches[0].screenY;
});

swipeZone.addEventListener("touchend", function (e) {
  touchEndX = e.changedTouches[0].screenX;
  touchEndY = e.changedTouches[0].screenY;

  const swipeX = Math.abs(touchEndX - touchStartX);
  const swipeY = Math.abs(touchEndY - touchStartY);

  let swipeDirection = "";

  if (swipeX > swipeY) {
    swipeDirection = touchEndX - touchStartX > 0 ? "right" : "left";
  } else {
    swipeDirection = touchEndY - touchStartY > 0 ? "down" : "up";
  }

  workerApi.postMessage({
    type: "swipe",
    direction: swipeDirection,
  });
});

//Utility functions
function isMobileViewport() {
  return window.matchMedia("(max-width: 767px)").matches;
}
