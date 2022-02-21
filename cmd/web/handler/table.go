package handler

import (
	"bigtable-viewer/package/db"
	"cloud.google.com/go/bigtable"
	"html/template"
	"log"
	"net/http"
)

func TableMultiplexer(c *db.Client, tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			showTable(c, tpl, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func showTable(client *db.Client, tpl *template.Template, w http.ResponseWriter, r *http.Request) {
	tableName := getFirstPathValue("/tables/(.+)", r.URL.Path)
	if tableName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err500 := func(err error) {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queryParams := r.URL.Query()
	start := getStringParam(queryParams, "start", "")
	search := getStringParam(queryParams, "search", "")
	limit := getIntParam(queryParams, "limit", 10)

	rows, err := client.FetchRows(tableName, start, limit, search)
	if err != nil {
		err500(err)
		return
	}

	structure, err := client.FetchFamilies(tableName)
	if err != nil {
		err500(err)
		return
	}

	lastRowKey := start
	if len(rows) > 0 {
		lastRowKey = rows[len(rows)-1].ID
	}

	err = printTemplate(w, tpl, "show-table.html", struct {
		TableName  string
		Families   []bigtable.FamilyInfo
		Rows       []db.Row
		LastRowKey string
		Limit      int
		Search     string
	}{
		TableName:  tableName,
		Families:   structure,
		Rows:       rows,
		LastRowKey: lastRowKey,
		Limit:      limit,
		Search:     search,
	})
	if err != nil {
		err500(err)
		return
	}
}
