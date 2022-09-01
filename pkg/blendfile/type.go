package blendfile

import (
	"fmt"
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
)

func (i *BlendFile) GetARGS() string {
	return fmt.Sprintf("%s %s", i.Exec.ARGS, i.Exec.ARGS)
}

func (i *BlendFile) GetAddonsAsARGS() string {
	addons := append(i.Exec.Addons, i.Addons...)
	return fmt.Sprintf("--addons=%s", strings.Join(addons[:], ","))
}

func (i *BlendFile) GetFullARGS() string {
	return fmt.Sprintf("%s %s %s", i.Path, i.GetARGS(), i.GetAddonsAsARGS())
}
