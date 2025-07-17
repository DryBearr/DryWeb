// ===============================================================
// File: engine.go
// Description: Combines interfaces together
// Author: DryBearr
// ===============================================================

package dryengine

// DryEngine is the main interface that combines rendering and event-handling capabilities.
// It embeds the DryRenderer and DryEvents interfaces to provide a unified engine abstraction.
type DryEngine interface {
	DryRenderer
	DryEvents
}
