package reference

import (
	"fmt"
	"path"
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
		// Basic check for now, we can add more checks later.
		if len(string(r)) <= len("local/") {
			return fmt.Errorf("invalid local reference: %s", r)
		}

		return nil
	}

	parts := strings.SplitN(string(r), "/", 4)
	if len(parts) < 4 {
		return fmt.Errorf("invalid reference: %s", r)
	}

	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("invalid reference: %s (contains empty parts)", r)
		}
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

// Parse parses a reference string and returns a Reference.
func Parse(s string) (Reference, error) {
	r := Reference(s)
	if err := r.Validate(); err != nil {
		return "", err
	}

	return r, nil
}

// Aliased resolves a reference string using a map of aliases.
// If the input string starts with an alias key, it expands the reference using the alias map.
func Aliased(input string, aliases map[string]string) (Reference, error) {
	for fullPath, alias := range aliases {
		if strings.HasPrefix(input, alias) {
			remainingPath := strings.TrimPrefix(input, alias)
			resolved := path.Join(fullPath, remainingPath)
			return Parse(resolved)
		}
	}

	return Parse(input)
}
