package blendfile

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/executable"
)

type (
	Client interface {
		FindExecutableByBuildReference(ref string) (*executable.Executable, error)
		FindAllAddonDirectories(ref []string) ([]string, error)
	}

	Config struct {
	}

	Service struct {
		conf *Config
		srv  Client
	}
)

func NewService(conf *Config, s Client) *Service {
	return &Service{
		conf: conf,
		srv:  s,
	}
}

func (s *Service) Load(path string) (*BlendFile, error) {
	ext := filepath.Ext(path)
	if ext != ".blend" {
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}

	c, err := os.ReadFile(filepath.Join(filepath.Dir(path), "rocketfile.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read rocketfile: %s", err)
	}

	var rkt RocketFile
	if err := json.Unmarshal(c, &rkt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	// Get build executable path.
	exec, err := s.srv.FindExecutableByBuildReference(rkt.Build)
	if err != nil {
		return nil, fmt.Errorf("failed to find executable: %s", err)
	}

	addons, err := s.srv.FindAllAddonDirectories(rkt.Packges)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addon directories: %s", err)
	}

	return &BlendFile{
		Exec:   exec,
		Path:   path,
		Addons: addons,
		ARGS:   rkt.ARGS,
	}, nil
}

func (s *Service) Save(file *BlendFile, safe bool) error {
	return fmt.Errorf("not implemented")
}

func (s *Service) Create(file *BlendFile) (*BlendFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Service) Open(file *BlendFile) error {
	fmt.Println(file.GetFullARGS())

	cmd := exec.Command(file.Exec.Path, file.GetFullARGS())

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open blend file: %s", err)
	}

	return nil
}
