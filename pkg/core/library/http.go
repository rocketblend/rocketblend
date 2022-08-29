package library

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type (
	LibraryConfig struct {
		UseProxy    bool
		BuildFile   string
		PackageFile string
		ReadMeFile  string
	}

	Client struct {
		conf *LibraryConfig
	}
)

func NewClient(conf *LibraryConfig) *Client {
	client := &Client{
		conf: conf,
	}

	return client
}

func NewDefaultConfig() *LibraryConfig {
	return &LibraryConfig{
		UseProxy:    false,
		BuildFile:   "build.json",
		PackageFile: "package.json",
		ReadMeFile:  "README.md",
	}
}

func (c *Client) FetchBuild(str string) (*Build, error) {
	rd, err := makeRequest(str, c.conf.BuildFile)
	if err != nil {
		return nil, err
	}

	var b *Build = &Build{}
	json.Unmarshal(rd, b)

	return b, nil
}

func (c *Client) FetchPackage(str string) (*Package, error) {
	rd, err := makeRequest(str, c.conf.PackageFile)
	if err != nil {
		return nil, err
	}

	var p *Package = &Package{}
	json.Unmarshal(rd, p)

	return p, nil
}

func makeRequest(str string, file string) ([]byte, error) {
	u, err := GetSourceUrl(str)
	if err != nil {
		return nil, err
	}

	u, err = url.JoinPath(u, file)
	if err != nil {
		return nil, err
	}

	r, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	rd, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return rd, nil
}
