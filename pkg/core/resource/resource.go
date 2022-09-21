package resource

import (
	_ "embed"
	"os"
	"path/filepath"
)

const Startup = "startup"

//go:embed resources/script.gopy
var startupScript string

type Resource struct {
	OutputPath string
	Content    string
}

type Service struct {
	Resources map[string]Resource
}

func NewService(dir string) *Service {
	return &Service{
		Resources: map[string]Resource{
			Startup: {
				OutputPath: filepath.Join(dir, Startup+".py"),
				Content:    startupScript,
			},
		},
	}
}

func (s *Service) SaveOut() error {
	for _, r := range s.Resources {
		file, err := os.Create(r.OutputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(r.Content)
		if err != nil {
			return err
		}
	}

	return nil
}
