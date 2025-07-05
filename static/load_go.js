/*
  ===============================================================
  File: load_go.js
  Description: loads wasm compiled in golang based on URL params
  Author: DryBearr
  ===============================================================
*/

import "./wasm_exec.js";

const DefaultWasm = "./wasm/game_of_life.wasm";
const AvailableWasms = new Map([
  ["game-of-life", "./wasm/game_of_life.wasm"],
  ["snake", "./wasm/snake.wasm"],
]);

let loadWasm = DefaultWasm;

const urlString = window.location.href;
const paramString = urlString.split("?")[1];
const queryString = new URLSearchParams(paramString);
let currentHref = "game-of-life";

if (queryString.size == 1 && queryString.has("wasm")) {
  const wasmParam = queryString.get("wasm");
  if (AvailableWasms.has(wasmParam)) {
    loadWasm = AvailableWasms.get(wasmParam);
    currentHref = wasmParam;
  }
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch(`${loadWasm}`), go.importObject)
  .then((result) => {
    go.run(result.instance);
  })
  .catch((error) => {
    //TODO: make this better
    console.error(error);
    setTimeout(() => {
      window.location.href = "/index.html";
    }, 5000);
  });

const nav = document.getElementById("wasm-links");
AvailableWasms.forEach((value, key) => {
  const link = document.createElement("a");
  link.innerText = `${key}`;
  link.href = `?wasm=${key}`;
  if (key === currentHref) {
    link.setAttribute("class", "active-link");
  }

  nav.append(link);
});
