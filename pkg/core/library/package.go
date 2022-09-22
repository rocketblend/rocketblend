package library

type (
	PackageSource struct {
		File string `json:"file"`
		URL  string `json:"url"`
	}

	Package struct {
		Reference string        `json:"reference"`
		Source    PackageSource `json:"source"`
	}
)
