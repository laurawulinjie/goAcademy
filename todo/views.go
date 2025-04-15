package main

import (
	"embed"
	"io/fs"
	"text/template"
)

var (
	//go:embed template/*.html
	templateFS embed.FS
	tmpl       *template.Template
)

type PageData struct {
	Todos map[int]Todo
}

func setupTemplate() error {
	subFS, err := fs.Sub(templateFS, "template")
	if err != nil {
		return err
	}
	tmpl, err = template.ParseFS(subFS, "*.html")
	return err
}
