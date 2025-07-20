/*
  ===============================================================
  File: main.ts
  Description: Entry point of web page scripts
  Author: DryBearr
  ===============================================================
*/

import {
  createGridBackground,
  handleBackgroundReSize,
  setBackgroundCanvas,
} from "./background";
import { setHeader, setHeaderNavWasmLinks } from "./header";
import { getCurrentActiveWasmLink } from "./util";
import "./index.css";

/*
  ===============================================================
  wasm modules list and root element setup
  ===============================================================
*/
const wasms: Array<string> = ["game_of_life", "snake"];
const wasmsToLoad = new Map<string, string>(
  wasms.map((name) => [name, `./wasm/${name}.wasm`]),
);

const root: HTMLElement | null = document.getElementById("root");

if (!root) throw new Error("no element with id `app` in the html document");

/*
  ===============================================================
  background canvas initialization
  ===============================================================
*/
const docWidth = document.documentElement.scrollWidth;
const docHeight = document.documentElement.scrollHeight;
const background = document.createElement("canvas");

setBackgroundCanvas(background);
createGridBackground(docWidth, docHeight);

root.append(background);

window.addEventListener("resize", () => {
  const docWidth = document.documentElement.scrollWidth;
  const docHeight = document.documentElement.scrollHeight;

  handleBackgroundReSize(docWidth, docHeight);
});

/*
  ===============================================================
  header and wasm navigation links
  ===============================================================
*/
const header: HTMLElement = document.createElement("header");

setHeader(header);
setHeaderNavWasmLinks(wasms);

root.append(header);

/*
  ===============================================================
  rendering canvas initialization 
  ===============================================================
*/

const canvasWidth = 800;
const canvasheight = 600;

const windowDiv = document.createElement("div");
windowDiv.setAttribute("class", "render-window");
root.append(windowDiv);

// Init Canvas and off screen canvas
const canvas = document.createElement("canvas");
canvas.setAttribute("id", "renderer");

canvas.setAttribute("width", canvasWidth.toString());
canvas.setAttribute("height", canvasheight.toString());
windowDiv.append(canvas);

/*
  ===============================================================
  worker initialization 
  ===============================================================
*/

// Start the worker api
let workerApi = new Worker("./worker_api.js", { type: "module" });

const activeWasm = getCurrentActiveWasmLink();

const loadWasm =
  wasmsToLoad.get(activeWasm ?? "") ?? "./wasm/game_of_life.wasm";

workerApi.postMessage({
  type: "init",
  wasm: loadWasm,
  width: canvasWidth,
  height: canvasheight,
});

// Initialize canvas rendering worker
const workerCanvas = new Worker("./worker_canvas.js");

const offScreenCanvas = canvas.transferControlToOffscreen();

workerCanvas.postMessage(
  {
    type: "init",
    offScreenCanvas: offScreenCanvas,
  },
  [offScreenCanvas],
);

/*
  ===============================================================
  events
  ===============================================================
*/

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

root.append(controlsDiv);

//Reload Wasm
const reloadWasmButton = document.createElement("button");
reloadWasmButton.textContent = "Reload";
reloadWasmButton.addEventListener("click", () => {
  workerApi.terminate();

  workerApi = new Worker("./worker_api.js", { type: "module" });

  workerApi.postMessage({
    type: "init",
    wasm: loadWasm,
    width: canvasWidth,
    height: canvasheight,
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
const getCanvasCoordinates = (event: MouseEvent | TouchEvent) => {
  const rect = canvas.getBoundingClientRect();
  if ("touches" in event && event.touches.length > 0) {
    return {
      x: Math.floor(event.touches[0].clientX - rect.left),
      y: Math.floor(event.touches[0].clientY - rect.top),
    };
  } else {
    return {
      x: Math.floor((event as MouseEvent).clientX - rect.left),
      y: Math.floor((event as MouseEvent).clientY - rect.top),
    };
  }
};

interface Point {
  x: number;
  y: number;
}

let isDraging = false;
let prevPoint: Point | null = null;

const handleDragStart = (event: MouseEvent | TouchEvent) => {
  event.preventDefault();

  isDraging = true;

  const { x, y } = getCanvasCoordinates(event);

  prevPoint = { x, y };
};

const handleDragMove = (event: MouseEvent | TouchEvent) => {
  event.preventDefault();
  if (!isDraging) return;

  const { x, y } = getCanvasCoordinates(event);

  if (prevPoint !== null) {
    workerApi.postMessage({
      type: "mouseDrag",
      x: prevPoint.x,
      y: prevPoint.y,
    });

    prevPoint = null;
  }

  workerApi.postMessage({ type: "mouseDrag", x, y });
};

const handleDragEnd = (event: MouseEvent | TouchEvent) => {
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

if (!swipeZone) throw new Error("can't find canvas for swipe zone");

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
