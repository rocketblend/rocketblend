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
	parts := strings.SplitN(string(r), "/", 4)
	if len(parts) < 4 {
		return fmt.Errorf("invalid reference: %s", r)
	}
	return nil
}

func (r Reference) Repo() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	parts := strings.SplitN(string(r), "/", 4)
	return strings.Join(parts[:3], "/"), nil
}

func (r Reference) RepoURL() (string, error) {
	repo, err := r.Repo()
	if err != nil {
		return "", err
	}

	return "https://" + repo, nil
}

func (r Reference) RepoPath() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	parts := strings.SplitN(string(r), "/", 4)
	return parts[3], nil
}

func Parse(s string) (Reference, error) {
	r := Reference(s)
	if err := r.Validate(); err != nil {
		return "", err
	}
	return r, nil
}
