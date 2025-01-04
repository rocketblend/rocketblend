package blender

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rocketblend/rocketblend/pkg/blender/parser"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (b *Blender) processOutput(output string) types.BlenderEvent {
	if output == "" {
		return nil
	}

	event, err := parser.ParseBlenderEvent(output)
	if err != nil {
		trimmedOutput := strings.ToLower(strings.TrimSpace(output))
		b.logger.Info("blender", map[string]interface{}{
			"output": trimmedOutput,
			"error":  err.Error(),
		})

		return &types.GenericEvent{Message: output}
	}

	eventMap := convertEventToMap(event)
	b.logger.Info("blender", eventMap)

	return event
}

func convertEventToMap(event types.BlenderEvent) map[string]interface{} {
	var result map[string]interface{}

	err := mapstructure.Decode(event, &result)
	if err != nil {
		return map[string]interface{}{
			"type":  "UnknownEvent",
			"error": err.Error(),
		}
	}

	switch event.(type) {
	case *types.RenderEvent:
		result["type"] = "render"
	case *types.GenericEvent:
		result["type"] = "generic"
	case *types.ErrorEvent:
		result["type"] = "error"
	default:
		result["type"] = "unknown"
	}

	return result
}
