/*
  ===============================================================
  File: background.ts
  Description: Provides utility function to draw background on canvas only one
  Author: DryBearr
  ===============================================================
*/

//TODO: In the future consider to use class, but for now this is enough

let background: HTMLCanvasElement | null = null;

const backgroundCellWidth = 70;
const backgroundCellHeight = 70;
const backgroundLineWidth = 1;
const backgroundColor = "#1E1E2E";
const backgroundLineColor = "#C0C0C0";

/**
 * Assigns the canvas element to be used for background rendering.
 * This must be called before using any drawing functions.
 *
 * @param canvas - The HTMLCanvasElement to use for background rendering.
 */
export function setBackgroundCanvas(canvas: HTMLCanvasElement) {
  background = canvas;
}

/**
 * Creates and draws a grid background on the canvas.
 * Sets canvas size, styles and draws grid aligned to center.
 *
 * @param width - The width of the canvas.
 * @param height - The height of the canvas.
 * @throws Error if the canvas has not been set via `setBackgroundCanvas`.
 */
export function createGridBackground(width: number, height: number) {
  if (!background) throw new Error("background canvas is not set yet");

  background.width = width;
  background.height = height;
  background.style.position = "absolute";
  background.style.top = "0";
  background.style.left = "0";
  background.style.zIndex = "-1000";
  background.style.margin = "0";
  background.style.padding = "0";
  background.style.pointerEvents = "none";
  background.style.overflow = "none";

  const ctx = background.getContext("2d");
  if (!ctx) throw new Error("failed to get canvas context");

  ctx.strokeStyle = backgroundLineColor;
  ctx.lineWidth = backgroundLineWidth;
  ctx.fillStyle = backgroundColor;

  ctx.fillRect(0, 0, width, height);

  const xOffset = (width % backgroundCellWidth) / 2;
  const yOffset = (height % backgroundCellHeight) / 2;

  // Vertical lines
  ctx.beginPath();
  for (let x = xOffset; x <= width; x += backgroundCellWidth) {
    const alignedX = Math.round(x) + 0.5;
    ctx.moveTo(alignedX, 0);
    ctx.lineTo(alignedX, height);
  }
  ctx.stroke();

  // Horizontal lines
  ctx.beginPath();
  for (let y = yOffset; y <= height; y += backgroundCellHeight) {
    const alignedY = Math.round(y) + 0.5;
    ctx.moveTo(0, alignedY);
    ctx.lineTo(width, alignedY);
  }
  ctx.stroke();
}

/**
 * Redraws the background grid after a resize event.
 * Can be used inside a resize handler to update the canvas dimensions and re-render the grid.
 *
 * @param width - New width for the canvas.
 * @param height - New height for the canvas.
 * @throws Error if the canvas has not been set via `setBackgroundCanvas`.
 */
export function handleBackgroundReSize(width: number, height: number) {
  if (!background) throw new Error("background canvas is not set yet");

  background.width = width;
  background.height = height;

  const ctx = background.getContext("2d");
  if (!ctx) throw new Error("failed to get canvas context");

  ctx.strokeStyle = backgroundLineColor;
  ctx.lineWidth = backgroundLineWidth;
  ctx.fillStyle = backgroundColor;
  ctx.fillRect(0, 0, width, height);

  const xOffset = (width % backgroundCellWidth) / 2;
  const yOffset = (height % backgroundCellHeight) / 2;

  // Vertical lines
  ctx.beginPath();
  for (let x = xOffset; x <= width; x += backgroundCellWidth) {
    const alignedX = Math.round(x) + 0.5;
    ctx.moveTo(alignedX, 0);
    ctx.lineTo(alignedX, height);
  }
  ctx.stroke();

  // Horizontal lines
  ctx.beginPath();
  for (let y = yOffset; y <= height; y += backgroundCellHeight) {
    const alignedY = Math.round(y) + 0.5;
    ctx.moveTo(0, alignedY);
    ctx.lineTo(width, alignedY);
  }
  ctx.stroke();
}
