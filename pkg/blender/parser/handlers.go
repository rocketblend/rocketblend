package parser

import (
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

// OperationHandler returns a specific BlenderEvent based on the operation details
type OperationHandler func(details string, base types.RenderBase) types.BlenderEvent

// Operation registry mapping operations to their handlers
var operationRegistry = map[string]OperationHandler{
	"rendering":     createRenderingEvent,
	"sample":        createRenderingEvent, // Treat "sample" as a rendering operation
	"synchronizing": createSynchronizingEvent,
	"updating":      createUpdatingEvent,
}

func createRenderingEvent(details string, base types.RenderBase) types.BlenderEvent {
	current, total := parseSamples(details)
	return &types.RenderingEvent{
		RenderBase: base,
		Current:    current,
		Total:      total,
		Operation:  "rendering",
	}
}

func createSynchronizingEvent(details string, base types.RenderBase) types.BlenderEvent {
	object := ""
	if strings.Contains(details, "object") {
		object = strings.TrimPrefix(details, "synchronizing object | ")
	}

	return &types.SynchronizingEvent{
		RenderBase: base,
		Object:     object,
	}
}

func createUpdatingEvent(details string, base types.RenderBase) types.BlenderEvent {
	return &types.UpdatingEvent{
		RenderBase: base,
		Details:    details,
	}
}
