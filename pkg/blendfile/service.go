package blendfile

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/install"
)

type (
	InstallService interface {
		FindByHash(hash string) (*install.Install, error)
	}

	Config struct {
		ARGS string
	}

	Service struct {
		conf Config
		srv  InstallService
	}
)

func NewService(conf Config) *Service {
	return &Service{
		conf: conf,
	}
}

func NewConfig(args string) Config {
	return Config{
		ARGS: args,
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

	var rkt *RocketFile
	if err := json.Unmarshal(c, &rkt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	inst, err := s.srv.FindByHash(rkt.Build)
	if err != nil {
		return nil, fmt.Errorf("failed to find build: %s", err)
	}

	return &BlendFile{
		Path:  path,
		Build: inst.Path,
		ARGS:  rkt.ARGS,
	}, nil
}

func (s *Service) Save(file *BlendFile, safe bool) error {
	return fmt.Errorf("not implemented")
}

func (s *Service) Create(file *BlendFile) (*BlendFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Service) Open(file *BlendFile) error {
	args := fmt.Sprintf("%s %s", file.ARGS, s.conf.ARGS)
	cmd := exec.Command(file.Build, args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open blend file: %s", err)
	}

	return nil
}
