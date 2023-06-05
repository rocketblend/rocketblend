package blendfile

import (
	"bufio"
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/renderoptions"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/runoptions"
)

//go:embed scripts/addonScript.gopy
var addonScript string

//go:embed scripts/createScript.gopy
var createScript string

type (
	Service interface {
		Render(ctx context.Context, blendFile *BlendFile, opts ...renderoptions.Option) error
		Run(ctx context.Context, blendFile *BlendFile, opts ...runoptions.Option) error
		Create(ctx context.Context, blendFile *BlendFile) error
	}

	Options struct {
		Logger       logger.Logger
		AddonScript  string
		CreateScript string
	}

	Option func(*Options)

	service struct {
		logger        logger.Logger
		addonScript   string
		createScript  string
		addonsEnabled bool
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func NewService(opts ...Option) (Service, error) {
	options := &Options{
		Logger:       logger.NoOp(),
		AddonScript:  addonScript,
		CreateScript: createScript,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.AddonScript == "" {
		return nil, fmt.Errorf("addon script is required")
	}

	if options.CreateScript == "" {
		return nil, fmt.Errorf("create script is required")
	}

	options.Logger.Debug("Initializing blendfile service")

	return &service{
		logger:       options.Logger,
		addonScript:  options.AddonScript,
		createScript: options.CreateScript,
	}, nil
}

func (s *service) Create(ctx context.Context, blendFile *BlendFile) error {
	err := os.MkdirAll(filepath.Dir(blendFile.FilePath), 0755)
	if err != nil {
		return err
	}

	script, err := parseOutputTemplate(s.createScript, map[string]string{
		"path": blendFile.FilePath,
	})
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, blendFile.Build.FilePath, "-b", "--python-expr", script)

	s.logger.Debug("running command", map[string]interface{}{"command": cmd.String()})

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			s.logger.Debug("Blender", map[string]interface{}{"Message": scanner.Text()})
		}
	}()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for command: %w", err)
	}

	err = Validate(blendFile)
	if err != nil {
		return fmt.Errorf("failed to validate: %s", err)
	}

	return nil
}

func (s *service) getCommand(ctx context.Context, blendFile *BlendFile, background bool, postArgs ...string) (*exec.Cmd, error) {
	// TODO: Only the rocketblend addon should be loaded. The addon will then load the other addons when blender starts. And will act as a toggle for addon support.
	preArgs := []string{}
	if background {
		preArgs = append(preArgs, "-b")
	}

	if blendFile.FilePath != "" {
		preArgs = append(preArgs, []string{blendFile.FilePath}...)
	}

	if s.addonsEnabled {
		json, err := json.Marshal(blendFile.Addons)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal addons: %s", err)
		}

		postArgs = append([]string{
			"--python-expr",
			s.addonScript,
		}, postArgs...)

		postArgs = append(postArgs, []string{
			"--",
			"-a",
			string(json),
		}...)
	}

	// Blender requires arguments to be in a specific order
	args := append(preArgs, postArgs...)
	cmd := exec.CommandContext(ctx, blendFile.Build.FilePath, args...)

	s.logger.Debug("running command", map[string]interface{}{"command": cmd.String()})

	return cmd, nil
}

func parseOutputTemplate(str string, data interface{}) (string, error) {
	// Define a new template with the input string
	tpl, err := template.New("").Parse(str)
	if err != nil {
		return "", err
	}

	// Execute the template with the data object and capture the output
	var buf strings.Builder
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
