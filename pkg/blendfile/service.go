package blendfile

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/core/executable"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
)

type (
	Client interface {
		FindResource(key string) (*resource.Resource, error)
		FindExecutableByBuildReference(ref string) (*executable.Executable, error)
		GetAddonMapByReferences(ref []string) (map[string]string, error)
	}

	Config struct {
		Debug bool
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

	addons, err := s.srv.GetAddonMapByReferences(rkt.Packages)
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
	script, err := s.srv.FindResource(resource.Startup)
	if err != nil {
		return fmt.Errorf("failed to find startup script: %s", err)
	}

	args := []string{
		file.Path,
		"--python",
		script.OutputPath,
	}

	addons := merge(file.Exec.Addons, file.Addons)
	json, err := json.Marshal(addons)
	if err != nil {
		return fmt.Errorf("failed to marshal addons: %s", err)
	}

	if len(addons) != 0 {
		args = append(args, []string{
			"--",
			"-a",
			string(json),
		}...)
	}

	cmd := exec.Command(file.Exec.Path, args...)

	if s.conf.Debug {
		fmt.Println(strings.ReplaceAll(cmd.String(), "\"", "\\\""))
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open blend file: %s", err)
	}

	return nil
}

func merge(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}
