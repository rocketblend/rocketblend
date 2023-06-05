package renderoptions

import "fmt"

type (
	Options struct {
		Background bool
		FrameStart int
		FrameEnd   int
		FrameStep  int
		Output     string
		Format     string
	}

	Option func(*Options)
)

func (opt *Options) Validate() error {
	if opt.FrameEnd < opt.FrameStart {
		return fmt.Errorf("invalid frame range: %d-%d:%d", opt.FrameStart, opt.FrameEnd, opt.FrameStep)
	}

	return nil
}

func WithBackground(background bool) Option {
	return func(options *Options) {
		options.Background = background
	}
}

func WithFrameRange(start, end, step int) Option {
	return func(options *Options) {
		options.FrameStart = start
		options.FrameEnd = end
		options.FrameStep = step
	}
}

func WithOutput(output string) Option {
	return func(options *Options) {
		options.Output = output
	}
}

func WithFormat(format string) Option {
	return func(options *Options) {
		options.Format = format
	}
}
