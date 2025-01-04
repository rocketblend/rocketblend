package types

type (
	// BlenderEvent represents a generic interface for all Blender events.
	BlenderEvent interface{}

	// GenericEvent represents a generic Blender event.
	GenericEvent struct {
		Message string `mapstructure:"message"`
	}

	// ErrorEvent represents an error Blender event.
	ErrorEvent struct {
		Message string `mapstructure:"message"`
	}

	// RenderBase represents common fields for all rendering-related Blender events.
	RenderBase struct {
		Frame      int    `mapstructure:"frame"`
		Memory     string `mapstructure:"memory"`
		PeakMemory string `mapstructure:"peakMemory"`
		Time       string `mapstructure:"time"`
	}

	// RenderingEvent represents a rendering-specific Blender event.
	RenderingEvent struct {
		RenderBase `mapstructure:",squash"`
		Current    int    `mapstructure:"current"`
		Total      int    `mapstructure:"total"`
		Operation  string `mapstructure:"operation"`
	}

	// SynchronizingEvent represents a synchronizing-specific Blender event.
	SynchronizingEvent struct {
		RenderBase `mapstructure:",squash"`
		Object     string `mapstructure:"object"`
	}

	// UpdatingEvent represents an updating-specific Blender event.
	UpdatingEvent struct {
		RenderBase
		Details string `mapstructure:"details"`
	}
)
