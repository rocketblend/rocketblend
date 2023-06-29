package rocketpack

import "github.com/rocketblend/rocketblend/pkg/downloader"

type Source struct {
	FileName string          `json:"fileName"`
	URI      *downloader.URI `json:"uri"`
}
