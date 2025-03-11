package cmd

import (
	"bytes"
	"text/template"
)

// buildModelfile builds a model file from the provided options.
func buildModelfile(opts runOptions) (string, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("modelfile").Parse(`FROM {{.Model}}
SYSTEM """{{.System}}"""
PARAMETER penalize_newline {{.Options.penalize_newline}}
PARAMETER seed {{.Options.seed}}
PARAMETER stop [{{range $index, $element := .Options.stop}}{{if $index}}, {{end}}{{$element}}{{end}}]
PARAMETER temperature {{.Options.temperature}}

{{range .Messages}}
MESSAGE {{.Role}} """{{.Content}}"""
{{end}}
`)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buf, opts)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
