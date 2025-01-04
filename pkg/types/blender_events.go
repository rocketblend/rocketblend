package types

type (
	// BlenderEvent represents a generic interface for all Blender events.
	BlenderEvent interface{}

	// GenericEvent represents a generic Blender event.
	GenericEvent struct {
		Message string `mapstructure:"message"`
	}

	ErrorEvent struct {
		Message string `mapstructure:"message"`
	}

	// RenderEvent represents a render-specific Blender event.
	RenderEvent struct {
		Frame      int               `mapstructure:"frame"`
		Memory     string            `mapstructure:"memory"`
		PeakMemory string            `mapstructure:"peakMemory"`
		Time       string            `mapstructure:"time"`
		Current    int               `mapstructure:"current"`
		Total      int               `mapstructure:"total"`
		Operation  string            `mapstructure:"operation"`
		Data       map[string]string `mapstructure:"data"`
	}
)
