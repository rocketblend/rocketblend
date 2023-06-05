package blendfile

import (
	_ "embed"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/renderoptions"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/runoptions"
	"github.com/rocketblend/rocketblend/pkg/driver/helpers"
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

	s.logger.Debug("Running command", map[string]interface{}{"command": cmd.String()})

	return cmd, nil
}

func (s *service) logAndReturnError(msg string, err error, fields ...map[string]interface{}) error {
	return helpers.LogAndReturnError(s.logger, msg, err, fields...)
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
