package install

type (
	Install struct {
		Id       string   `json:"id"`
		Build    string   `json:"build"`
		Path     string   `json:"path"`
		Packages []string `json:"packages"`
		CheckSum string   `json:"checksum"`
	}
)
