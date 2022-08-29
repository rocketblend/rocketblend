package library

type (
	Source struct {
		Platform   string `json:"platform"`
		Executable string `json:"executable"`
		URL        string `json:"url"`
	}

	Build struct {
		Args     string   `json:"args"`
		Source   []Source `json:"source"`
		Packages []string `json:"packages"`
	}
)

type (
	Package struct {
		Source string `json:"source"`
	}
)

type (
	Install struct {
		Path string `json:"path"`
	}

	Pack struct {
		Path string `json:"path"`
	}
)
