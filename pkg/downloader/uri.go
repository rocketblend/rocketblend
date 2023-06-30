package downloader

import (
	"encoding/json"
	"net/url"
)

type URI url.URL

func (u *URI) IsLocal() bool {
	return u.Scheme == "file"
}

func (u *URI) IsRemote() bool {
	return u.Scheme == "http" || u.Scheme == "https"
}

func (u *URI) String() string {
	return (*url.URL)(u).String()
}

func (u *URI) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URI) UnmarshalJSON(data []byte) error {
	var urlString string
	err := json.Unmarshal(data, &urlString)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	*u = URI(*parsedURL)
	return nil
}

func NewURI(uri string) (*URI, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return (*URI)(u), nil
}
