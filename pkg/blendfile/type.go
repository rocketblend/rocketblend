package blendfile

import (
	"fmt"
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

func (i *BlendFile) GetARGS() string {
	return fmt.Sprintf("%s %s", i.Exec.ARGS, i.Exec.ARGS)
}

func (i *BlendFile) GetAddonStr() string {
	addons := append(i.Exec.Addons, i.Addons...)
	return strings.Join(addons[:], ",")
}

func (i *BlendFile) GetPythonArgs() string {
	dirname, _ := os.UserHomeDir()
	script := filepath.Join(dirname, ".rocketblend", "scripts", "arg_script.py")
	return fmt.Sprintf("--python \"%s\" -- -a \"%s\"", script, i.GetAddonStr())
}
