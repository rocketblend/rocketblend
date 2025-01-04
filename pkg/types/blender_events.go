package types

type (
	// BlenderEvent represents a generic interface for all Blender events.
	BlenderEvent interface{}

	// GenericEvent represents a generic Blender event.
	GenericEvent struct {
		Message string
	}

	// RenderEvent represents a render-specific Blender event.
	RenderEvent struct {
		Frame      int
		Memory     string
		PeakMemory string
		Time       string
		Operation  string
		Data       map[string]string
	}
)
