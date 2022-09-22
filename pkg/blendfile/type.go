package blendfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/core/executable"
)

type (
	RocketFile struct {
		Build   string
		ARGS    string
		Version string
		Packges []string
	}

	BlendFile struct {
		Exec   *executable.Executable
		Path   string
		Addons []string
		ARGS   string
	}

	AddonDict struct {
		Name string
		Path string
	}
)

func (i *BlendFile) Get() []string {
	args := []string{i.Path}

	a := append(i.Exec.Addons, i.Addons...)
	addons := strings.Join(a[:], ",")

	if addons != "" {
		d, _ := os.UserHomeDir()
		script := filepath.Join(d, ".rocketblend", "scripts", "arg_script.py")

		args = append(args, []string{
			"--python",
			script,
			"--",
			"-a",
			addons,
		}...)
	}

	return args
}
