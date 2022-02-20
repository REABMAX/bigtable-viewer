package handler

import (
	"bytes"
	"html/template"
	"net/http"
)

func printTemplate(w http.ResponseWriter, tpl *template.Template, name string, data interface{}) error {
	buf := &bytes.Buffer{}
	err := tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
