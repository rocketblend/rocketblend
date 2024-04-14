package blender

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	renderArguments struct {
		start   int
		end     int
		step    int
		output  string
		format  types.RenderFormat
		devices []types.CyclesDevice
		threads int
	}

	rocketblendArguments struct {
		addons []*types.Installation
	}

	arguments struct {
		background    bool
		blendFilePath string
		script        string
		render        *renderArguments
		rockeblend    *rocketblendArguments
	}
)

func (a *renderArguments) ARGS() []string {
	if a.start == 0 && a.end == 0 {
		return nil
	}

	args := []string{}
	if a.start != 0 {
		args = append(args, "--frame-start", fmt.Sprint(a.start))
	}

	if a.end != 0 {
		args = append(args, "--frame-end", fmt.Sprint(a.end))
	}

	if a.step != 0 {
		args = append(args, "--frame-jump", fmt.Sprint(a.step))
	}

	if a.output != "" {
		args = append(args, "--render-output", a.output)
	}

	if a.format != "" {
		args = append(args, "--render-format", string(a.format), "-x", "1")
	}

	if len(a.devices) > 0 {
		devices := []string{"--cycles-device"}
		for _, device := range a.devices {
			devices = append(devices, string(device))
		}

		args = append(args, strings.Join(devices, "+"))
	}

	if a.threads > 0 {
		args = append(args, "-t", strconv.Itoa(a.threads))
	}

	return append(args, "-a")
}

func (a *rocketblendArguments) ARGS() []string {
	if a.addons != nil {
		json, err := json.Marshal(a.addons)
		if err != nil {
			return nil
		}

		return []string{
			"--",
			"-a",
			string(json),
		}
	}

	return nil
}

func (a *arguments) ARGS() []string {
	args := []string{}
	if a.background {
		args = append(args, "-b")
	}

	if a.blendFilePath != "" {
		args = append(args, a.blendFilePath)
	}

	if a.script != "" {
		args = append(args, []string{
			"--python-expr",
			a.script,
		}...)
	}

	if a.render != nil {
		args = append(args, a.render.ARGS()...)
	}

	if a.script != "" && a.rockeblend != nil {
		args = append(args, a.rockeblend.ARGS()...)
	}

	return args
}
