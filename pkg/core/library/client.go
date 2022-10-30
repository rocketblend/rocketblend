package library

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type (
	LibraryConfig struct {
		UseProxy bool
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

func NewClientConfig() *LibraryConfig {
	return &LibraryConfig{
		UseProxy: false,
	}
}

func (c *Client) FetchBuild(str string) (*Build, error) {
	rd, err := c.makeRequest(str, BuildFile)
	if err != nil {
		return nil, err
	}

	var b *Build = &Build{}
	json.Unmarshal(rd, b)

	return b, nil
}

func (c *Client) FetchPackage(str string) (*Package, error) {
	rd, err := c.makeRequest(str, PackgeFile)
	if err != nil {
		return nil, err
	}

	var p *Package = &Package{}
	json.Unmarshal(rd, p)

	return p, nil
}

func (c *Client) makeRequest(str string, file string) ([]byte, error) {
	u, err := GetBuildUrl(str)
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
