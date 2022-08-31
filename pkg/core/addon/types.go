package addon

type (
	Addon struct {
		Id       string `json:"id"`
		Package  string `json:"package"`
		Path     string `json:"path"`
		CheckSum string `json:"checksum"`
	}
)
