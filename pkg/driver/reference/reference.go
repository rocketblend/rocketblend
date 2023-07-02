package reference

import (
	"fmt"
	"strings"
)

type Reference string

func (r *Reference) String() string {
	return string(*r)
}

func (r Reference) IsLocalOnly() bool {
	return strings.HasPrefix(string(r), "local")
}

func (r Reference) Validate() error {
	if r.IsLocalOnly() {
		return nil
	}

	parts := strings.SplitN(string(r), "/", 4)
	if len(parts) < 4 {
		return fmt.Errorf("invalid reference: %s", r)
	}

	return nil
}

func (r Reference) GetRepo() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	if r.IsLocalOnly() {
		return "local/", nil
	}

	parts := strings.SplitN(string(r), "/", 4)
	return strings.Join(parts[:3], "/"), nil
}

func (r Reference) GetRepoURL() (string, error) {
	repo, err := r.GetRepo()
	if err != nil {
		return "", err
	}

	if r.IsLocalOnly() {
		return "", fmt.Errorf("local reference: %s has no URL", r)
	}

	return "https://" + repo, nil
}

func (r Reference) GetRepoPath() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	if r.IsLocalOnly() {
		// For local references, just return the entire reference as the path
		return string(r)[len("local"):], nil // remove the 'local' from the beginning
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
