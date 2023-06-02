package reference

import (
	"fmt"
	"strings"
)

type Reference string

func (r Reference) String() string {
	return string(r)
}

func (r Reference) Validate() error {
	parts := strings.SplitN(string(r), "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid reference: %s", r)
	}
	return nil
}

func (r Reference) RepoURL() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	parts := strings.SplitN(string(r), "/", 2)
	return parts[0], nil
}

func (r Reference) RepoPath() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	parts := strings.SplitN(string(r), "/", 2)
	return parts[1], nil
}
