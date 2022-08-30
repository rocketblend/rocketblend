package blendfile

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
)

type (
	Client interface {
		FindInstall(build string) (*install.Install, error)
		FindBuildByPath(build string) (*library.Build, error)
		Platform() runtime.Platform
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

	// Check if the build is installed
	inst, err := s.srv.FindInstall(rkt.Build)
	if err != nil {
		return nil, fmt.Errorf("failed to find install: %s", err)
	}

	// Get build information from local install
	build, err := s.srv.FindBuildByPath(inst.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to find build information: %s", err)
	}

	return &BlendFile{
		Path:  path,
		Build: filepath.Join(inst.Path, build.GetSourceForPlatform(s.srv.Platform()).Executable),
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
	args := fmt.Sprintf("%s %s", file.Path, file.ARGS)
	cmd := exec.Command(file.Build, args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open blend file: %s", err)
	}

	return nil
}
