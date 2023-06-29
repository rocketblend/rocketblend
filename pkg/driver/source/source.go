package source

import "net/url"

type Source struct {
	FileName string   `json:"fileName"`
	URI      *url.URL `json:"uri"`
}

func (s *Source) IsLocal() bool {
	return s.URI.Scheme == "file"
}

func (s *Source) IsRemote() bool {
	return s.URI.Scheme == "http" || s.URI.Scheme == "https"
}

func (s *Source) String() string {
	return s.URI.String()
}
