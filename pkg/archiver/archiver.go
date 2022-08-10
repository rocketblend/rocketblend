package archiver

import (
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

func Extract(filePath string, deleteArchive bool) error {
	fileDir := filepath.Dir(filePath)
	err := archiver.Unarchive(filePath, fileDir)
	if err != nil {
		return err
	}

	if deleteArchive {
		err = os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
