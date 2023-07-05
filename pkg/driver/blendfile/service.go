package blendfile

import (
	_ "embed"
	"strings"

	"context"
	"fmt"
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
		Logger        logger.Logger
		addonsEnabled bool
		AddonScript   string
		CreateScript  string
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

func WithAddonsEnabled(enabled bool) Option {
	return func(o *Options) {
		o.addonsEnabled = enabled
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

	options.Logger.Debug("Initializing blendfile service", map[string]interface{}{"addonsEnabled": options.addonsEnabled})

	return &service{
		logger:        options.Logger,
		addonsEnabled: options.addonsEnabled,
		addonScript:   options.AddonScript,
		createScript:  options.CreateScript,
	}, nil
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
