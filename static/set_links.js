/*
  ===============================================================
  File: load_go.js
  Description: Set wasm navigation links
  Author: DryBearr
  ===============================================================
*/

import config from "./available_wasms.json" with { type: "json" };

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
