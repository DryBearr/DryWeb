// ===============================================================
// File: colors.go
// Description: Provide utility function for colors
// Author: DryBearr
// ===============================================================

package util

import (
	"encoding/hex"
	"wasm/dryeve/models"
)

// EncodeColorHex returns the RGBA color as a hex string in the format #rrggbbaa.
func EncodeColorHex(p models.Pixel) string {
	b := []byte{p.R, p.G, p.B, p.A}
	return "#" + hex.EncodeToString(b)
}
