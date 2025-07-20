/*
  ===============================================================
  File: header.ts
  Description: Defines web page header and provides function to manipulate it 
  Author: DryBearr
  ===============================================================
*/

import { getCurrentActiveWasmLink } from "./util";

let header: HTMLElement | null = null;

/**
 * Sets the reference to the document's <head> element.
 * This must be called before adding navigation links.
 *
 * @param headerElement - The HTML head element to be used for injecting navigation.
 */
export function setHeader(headerElement: HTMLElement) {
  header = headerElement;
}

/**
 * Creates a navigation bar in the <head> that allows switching between WASM modules.
 * If the current active WASM link is not specified in the URL or not found in the list,
 * the first link will be used as the default.
 *
 * @param wasmLinks - An array of string names of WASM modules to be linked in navigation.
 * @throws Will throw an error if the array is empty or undefined.
 */
export function setHeaderNavWasmLinks(wasmLinks: Array<string>) {
  if (!wasmLinks || wasmLinks.length === 0) {
    throw new Error("no wasm links provided for navigation");
  }

  createHeaderNav(wasmLinks);
}

/**
 * Internal helper to create the navigation bar with WASM module links.
 * Clears existing navigation and appends new links.
 *
 * @param wasmLinks - List of wasm module names to generate links for.
 */
function createHeaderNav(wasmLinks: Array<string>) {
  if (!header)
    throw new Error("can't add wasm links: head element does not exist");

  let nav = document.getElementById("header-wasm-nav");
  if (!nav) {
    nav = document.createElement("nav");
    nav.setAttribute("id", "header-wasm-nav");
    header.append(nav);
  } else {
    nav.innerHTML = "";
  }

  let currentRef = getCurrentActiveWasmLink();
  if (currentRef === "") {
    currentRef = wasmLinks[0];
  } else if (wasmLinks.filter((value) => value === currentRef).length == 0) {
    currentRef = wasmLinks[0];
  }

  wasmLinks.forEach((val) => {
    const link = document.createElement("a");
    const span = document.createElement("span");

    span.textContent = val;
    link.href = `?wasm=${val}`;
    link.appendChild(span);

    if (val === currentRef) {
      link.classList.add("active-link");
    }

    nav.appendChild(link);
  });
}
