package rocketpack

import "github.com/rocketblend/rocketblend/pkg/downloader"

type Source struct {
	Resource string          `json:"resource,omitempty"`
	URI      *downloader.URI `json:"uri,omitempty"`
}
