package types

type (
	Container interface {
		GetLogger() (Logger, error)
		GetValidator() (Validator, error)
		GetDownloader() (Downloader, error)
		GetExtractor() (Extractor, error)
		GetConfigurator() (Configurator, error)
		GetRepository() (Repository, error)
		GetBlender() (Blender, error)
	}
)