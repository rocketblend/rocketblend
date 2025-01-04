package parser

import (
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

// OperationHandler is a function that processes specific operation details
type OperationHandler func(details string, event *types.RenderEvent)

// Operation registry mapping operations to their handlers
var operationRegistry = map[string]OperationHandler{
	"rendering":     handleRendering,
	"sample":        handleRendering, // Treat "sample" as a rendering operation
	"synchronizing": handleSynchronizing,
	"initializing":  handleGenericOperation,
	"waiting":       handleGenericOperation,
	"updating":      handleUpdating,
	"loading":       handleGenericOperation,
}

// Handles the specified operation based on the registry
func handleOperation(operationRaw string, event *types.RenderEvent) {
	for op, handler := range operationRegistry {
		if strings.Contains(operationRaw, op) {
			event.Operation = op
			handler(operationRaw, event)
			return
		}
	}

	event.Operation = "general"
	handleGenericOperation(operationRaw, event)
}

// Specific handlers for each operation type
func handleRendering(details string, event *types.RenderEvent) {
	currentSample, totalSamples := parseSamples(details)
	event.Current = currentSample
	event.Total = totalSamples
}

func handleSynchronizing(details string, event *types.RenderEvent) {
	if strings.Contains(details, "object") {
		event.Data["object"] = strings.TrimPrefix(details, "synchronizing object | ")
	}
}

func handleUpdating(details string, event *types.RenderEvent) {
	event.Data["details"] = details
}

func handleGenericOperation(details string, event *types.RenderEvent) {
	event.Data["details"] = details
}
