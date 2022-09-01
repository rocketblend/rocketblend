package library

type (
	Package struct {
		Reference string `json:"reference"`
		Name      string `json:"name"`
		Source    string `json:"source"`
	}
)
