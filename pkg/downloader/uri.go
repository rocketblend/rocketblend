package downloader

import "net/url"

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

func NewURI(uri string) (*URI, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return (*URI)(u), nil
}
