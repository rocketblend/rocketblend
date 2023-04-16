package resource

import (
	_ "embed"
	"html/template"
	"strings"
)

//go:embed resources/addonScript.gopy
var addonScript string

//go:embed resources/createScript.gopy
var createScript string

type Service struct {
	addonScript  string
	createScript string
}

func NewService() *Service {
	return &Service{
		addonScript:  addonScript,
		createScript: createScript,
	}
}

func (s *Service) GetAddonScript() string {
	return s.addonScript
}

func (s *Service) GetCreateScript(path string) (string, error) {
	vars := map[string]string{
		"path": path,
	}

	// Render the template
	output, err := renderTemplate(s.createScript, vars)
	if err != nil {
		return "", err
	}

	return output, nil
}

// renderTemplate takes a template string and a map of variables as input
// and returns the output as a string
func renderTemplate(tmplStr string, vars map[string]string) (string, error) {
	// Create a new template with a name
	t := template.Must(template.New("template").Parse(tmplStr))

	// Create a buffer to write the output
	var b strings.Builder

	// Execute the template and inject the variables
	if err := t.ExecuteTemplate(&b, "tmpl", vars); err != nil {
		return "", err
	}

	// Return the output
	return b.String(), nil
}
