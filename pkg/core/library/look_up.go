package library

import (
	"fmt"
	"net/url"
	"strings"
)

const hostname string = "github.com"
const rawHostname string = "raw.githubusercontent.com"

type hostError struct {
	Host string
}

func (e *hostError) Error() string {
	return fmt.Sprintf("invalid host: %s", e.Host)
}

func GetBuildUrl(str string) (string, error) {
	str = addProtcol(str)

	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	if err := validateHost(u); err != nil {
		return "", err
	}

	u.Host = rawHostname
	u.Path = addToPathIndex(u.Path, 3, "master")

	return u.String(), nil
}

func addAtIndex(s []string, index int, value string) []string {
	return append(s[:index], append([]string{value}, s[index:]...)...)
}

func addToPathIndex(path string, index int, value string) string {
	s := strings.Split(path, "/")
	s = addAtIndex(s, index, value)
	return strings.Join(s, "/")
}

func addProtcol(str string) string {
	if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
		str = "https://" + str
	}

	return str
}

func validateHost(u *url.URL) error {
	if u.Host != hostname {
		return &hostError{Host: u.Host}
	}

	return nil
}
