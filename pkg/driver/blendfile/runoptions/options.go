package runoptions

type (
	Options struct {
		Background bool
	}

	Option func(*Options)
)

func WithBackground(background bool) Option {
	return func(options *Options) {
		options.Background = background
	}
}

func (opt *Options) Validate() error {
	return nil
}
