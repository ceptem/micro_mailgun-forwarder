package main

import (
	"fmt"
	"os"
	"text/template"
)

func LoadTemplate(name string, body string) *template.Template {
	template, err := template.New(name).Parse(body)
	if err != nil {
		fmt.Println("Alert: Failed to parse template:", err)
		os.Exit(1)
	}
	return template
}
