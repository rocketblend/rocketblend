package blender

import "github.com/rocketblend/rocketblend/pkg/types"

const (
	Default   RenderEngine = ""
	Eevee     RenderEngine = "BLENDER_EEVEE"
	Workbench RenderEngine = "BLENDER_WORKBENCH"
	Cycles    RenderEngine = "CYCLES"
)

type (
	RenderEngine string
)

func convertRenderEngine(engine types.RenderEngine) RenderEngine {
	switch engine {
	case types.Eevee:
		return Eevee
	case types.Workbench:
		return Workbench
	case types.Cycles:
		return Cycles
	default:
		return Default
	}
}
