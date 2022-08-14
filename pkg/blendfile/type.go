package blendfile

type (
	RocketFile struct {
		Build   string
		ARGS    string
		Version string
	}

	BlendFile struct {
		Path  string
		Build string
		ARGS  string
	}
)
