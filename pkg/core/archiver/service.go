package archiver

import (
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

type (
	Config struct {
		DeleteArchive bool
	}

	Service struct {
		conf Config
	}
)

func NewService(conf *Config) *Service {
	return &Service{
		conf: *conf,
	}
}

func NewConfig(deleteArchive bool) *Config {
	return &Config{
		DeleteArchive: deleteArchive,
	}
}

func (s *Service) Extract(path string) error {
	dir := filepath.Dir(path)
	err := archiver.Unarchive(path, dir)
	if err != nil {
		return err
	}

	if s.conf.DeleteArchive {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
