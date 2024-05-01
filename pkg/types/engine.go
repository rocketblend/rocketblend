package types

type (
	RenderEngine string
)

const (
	Default   RenderEngine = ""
	Eevee     RenderEngine = "eevee"
	Workbench RenderEngine = "workbench"
	Cycles    RenderEngine = "cycles"
)
