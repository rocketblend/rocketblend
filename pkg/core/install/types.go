package install

type (
	Install struct {
		Id       string `json:"id"`
		Build    string `json:"build"`
		Path     string `json:"path"`
		CheckSum string `json:"checksum"`
	}

	Pack struct {
		Id       string `json:"id"`
		Path     string `json:"path"`
		Package  string `json:"package"`
		CheckSum string `json:"checksum"`
	}
)

// func (i *Install) GetExecutableForPlatform(platform string) string {
// 	return filepath.Join(i.Path, i.Build.GetSourceForPlatform(platform).Executable)
// }
