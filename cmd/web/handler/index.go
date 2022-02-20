package handler

import (
	"bigtable-viewer/package/db"
	"html/template"
	"log"
	"net/http"
)

func IndexMultiplexer(c *db.Client, tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			showIndex(c, tpl, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func showIndex(client *db.Client, tpl *template.Template, w http.ResponseWriter, r *http.Request) {
	tables, err := client.FetchTables()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = printTemplate(w, tpl, "index.html", struct {
		Tables []string
	}{
		Tables: tables,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}