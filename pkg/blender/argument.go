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
		Start   int
		End     int
		Step    int
		Output  string
		Format  RenderFormat
		Devices []CyclesDevice
		Engine  RenderEngine
		Threads int
	}

	rocketblendArguments struct {
		Addons []*types.Installation
		Strict bool
	}

	arguments struct {
		Background    bool
		BlendFilePath string
		Script        string
		Render        *renderArguments
		Rockeblend    *rocketblendArguments
	}
)

func (a *renderArguments) ARGS() []string {
	if a.Start == 0 && a.End == 0 {
		return nil
	}

	args := []string{}
	if a.Engine != "" {
		args = append(args, "--engine", string(a.Engine))
	}

	if a.Start != 0 {
		args = append(args, "--frame-start", fmt.Sprint(a.Start))
	}

	if a.End != 0 {
		args = append(args, "--frame-end", fmt.Sprint(a.End))
	}

	if a.Step != 0 {
		args = append(args, "--frame-jump", fmt.Sprint(a.Step))
	}

	if a.Output != "" {
		args = append(args, "--render-output", a.Output)
	}

	if a.Format != "" {
		args = append(args, "--render-format", string(a.Format), "-x", "1")
	}

	if len(a.Devices) > 0 {
		devices := []string{"--cycles-device"}
		for _, device := range a.Devices {
			devices = append(devices, string(device))
		}

		args = append(args, strings.Join(devices, "+"))
	}

	if a.Threads > 0 {
		args = append(args, "-t", strconv.Itoa(a.Threads))
	}

	return append(args, "-a")
}

func (a *rocketblendArguments) ARGS() []string {
	args := []string{}
	if a.Addons != nil {
		json, err := json.Marshal(a.Addons)
		if err != nil {
			return nil
		}

		args = append(args, []string{
			"-a",
			string(json),
		}...)
	}

	if a.Strict {
		args = append(args, []string{
			"-s",
		}...)
	}

	return args
}

func (a *arguments) ARGS() []string {
	args := []string{}
	if a.Background {
		args = append(args, "-b")
	}

	if a.BlendFilePath != "" {
		args = append(args, a.BlendFilePath)
	}

	if a.Script != "" {
		args = append(args, []string{
			"--python-expr",
			a.Script,
		}...)
	}

	if a.Render != nil {
		args = append(args, a.Render.ARGS()...)
	}

	if a.Script != "" && a.Rockeblend != nil {
		args = append(args, "--")
		args = append(args, a.Rockeblend.ARGS()...)
	}

	return args
}
