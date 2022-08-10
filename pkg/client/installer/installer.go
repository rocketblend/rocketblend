package installer

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/archiver"
	"github.com/rocketblend/rocketblend/pkg/downloader"
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
