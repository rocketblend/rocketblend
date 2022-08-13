package build

import (
	"time"

	"github.com/rocketblend/rocketblend/pkg/core/remote"
)

type (
	FetchRequest struct {
		Remotes  []*remote.Remote
		Platform string
		Version  string
		Tag      string
	}
)

type (
	Build struct {
		Platform    string    `json:"platform"`
		Name        string    `json:"name"`
		Version     string    `json:"version"`
		Tag         string    `json:"tag"`
		Hash        string    `json:"hash"`
		DownloadUrl string    `json:"downloadurl"`
		CrawledAt   time.Time `json:"crawledat"`
	}

	Response struct {
		Data []Build `json:"data"`
	}
)
