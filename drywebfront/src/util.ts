/*
  ===============================================================
  File: util.ts
  Description: Provides utility functions used across web page scripts 
  Author: DryBearr
  ===============================================================
*/

/**
 * Extracts the current active WebAssembly (wasm) module name
 * from the URL's `?wasm=...` query parameter.
 *
 * @returns The value of the `wasm` parameter in the URL, or an empty string if not set.
 */
export function getCurrentActiveWasmLink(): string {
  const urlParams = new URLSearchParams(window.location.search);

  return urlParams.get("wasm") || "";
}
