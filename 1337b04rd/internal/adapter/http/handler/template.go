package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	tmplPath := filepath.Join("web", tmpl) // Это создаст правильный путь для шаблона
	fmt.Println(tmplPath)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}
