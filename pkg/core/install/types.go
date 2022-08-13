package install

type (
	FindRequest struct {
		SortBy []string
	}
)

type (
	Install struct {
		Path    string `json:"path"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Hash    string `json:"hash"`
	}
)
