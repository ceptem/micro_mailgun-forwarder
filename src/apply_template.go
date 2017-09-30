package main

import (
	"bytes"
	"text/template"
)

func ApplyTemplate(
	tpl *template.Template,
	context Context,
) (string, error) {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, context)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
