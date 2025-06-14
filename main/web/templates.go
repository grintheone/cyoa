package main

import (
	"html/template"
	"net/http"
	"path"
)

func isEmpty(slice []Option) bool {
	return len(slice) == 0
}

var functions = template.FuncMap{
	"isEmpty": isEmpty,
}

func renderTemplate(w http.ResponseWriter, s Story) error {
	p := path.Join("ui", "html", "base.html")

	tmpl := template.New("").Funcs(functions)
	tmpl, err := tmpl.ParseFiles(p)
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, "base.html", s)
	if err != nil {
		return err
	}

	return nil
}
