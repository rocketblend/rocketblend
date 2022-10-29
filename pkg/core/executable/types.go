package executable

type (
	Executable struct {
		Path   string
		Addons map[string]string
		ARGS   string
	}
)
