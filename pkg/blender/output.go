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
		b.logger.Debug("blender", map[string]interface{}{
			"output": trimmedOutput,
			"error":  err.Error(),
		})

		return err
	}

	eventMap := convertEventToMap(event)
	if len(eventMap) != 0 {
		b.logger.Info("blender", eventMap)
	}

	return event
}

func convertEventToMap(event types.BlenderEvent) map[string]interface{} {
	result := make(map[string]interface{})

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &result,
		TagName:  "mapstructure",
		Squash:   true,
	})
	if err != nil {
		return map[string]interface{}{
			"type":  "unknown",
			"error": "failed to create decoder",
		}
	}

	err = decoder.Decode(event)
	if err != nil {
		return map[string]interface{}{
			"type":  "unknown",
			"error": err.Error(),
		}
	}

	switch event.(type) {
	case *types.QuitEvent:
		result["event"] = "quit"
	case *types.SavedFileEvent:
		result["event"] = "saved"
	case *types.RenderingEvent:
		result["event"] = "rendering"
	case *types.SynchronizingEvent:
		result["event"] = "synchronizing"
	case *types.UpdatingEvent:
		result["event"] = "updating"
	case *types.GenericEvent:
		result["event"] = "raw"
	case *types.ErrorEvent:
		result["event"] = "error"
	}

	return result
}
