package installer

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/vendor/github.com/mholt/archiver/v3"
)

func Install(fileUrl string, downloadDir string) error {
	fmt.Println("Download Started")

	err := downloader.DownloadFile(downloadDir, fileUrl)
	if err != nil {
		return err
	}

	fmt.Println("Download Finished")

	fmt.Println("Extract Started")
	err = archiver.Extract(downloadDir, true)
	if err != nil {
		return err
	}

	fmt.Println("Extract Finished")

	return nil
}
