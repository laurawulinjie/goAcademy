package main

import (
	"embed"
	"io/fs"
	"text/template"
)

var (
	//go:embed pages/static
	staticFS embed.FS
	//go:embed pages/dynamic
	dynamicFS embed.FS
	tmpl      *template.Template
)

type PageData struct {
	Todos map[int]Todo
}

func getStaticFS() (fs.FS, error) {
	return fs.Sub(staticFS, "pages/static")
}

func getDynamicFS() (fs.FS, error) {
	return fs.Sub(dynamicFS, "pages/dynamic")
}

func setupDynamicPages() error {
	subFS, err := getDynamicFS()
	if err != nil {
		return err
	}
	tmpl, err = template.ParseFS(subFS, "*.html")
	return err
}
