package library

type (
	PackageSource struct {
		File string `json:"file"`
		URL  string `json:"url"`
	}

	Package struct {
		Reference string        `json:"reference"`
		Name      string        `json:"name"`
		Source    PackageSource `json:"source"`
	}
)
