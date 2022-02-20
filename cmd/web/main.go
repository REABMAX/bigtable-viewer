package main

import (
	"bigtable-viewer/cmd/web/handler"
	"bigtable-viewer/internal/settings"
	"bigtable-viewer/package/db"
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates
var templatesDir embed.FS

//go:embed assets
var assetsDir embed.FS

func main() {
	mux := http.NewServeMux()
	ctx := context.Background()
	project, instance, errors := settings.All()
	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		log.Fatalln("missing or wrong arguments: see above")
	}

	client, err := db.NewClient(ctx, project, instance)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	tpl := template.Must(template.ParseFS(templatesDir, "templates/*.html"))
	mux.Handle("/assets", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsDir))))
	mux.HandleFunc("/tables/", handler.TableMultiplexer(client, tpl))
	mux.HandleFunc("/", handler.IndexMultiplexer(client, tpl))

	log.Fatalln(http.ListenAndServe(":3001", mux))
}
