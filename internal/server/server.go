package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/internal/database"
	handler2 "github.com/vigneshrajj/gofind/internal/handlers"
)

func StartServer() {
	_, db, err := database.NewDBConnection(config.DbPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = database.EnsureDefaultCommandsExist(db)
	if err != nil {
		log.Fatal(err)
		return
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		handler2.HandleQuery(w, r, query, db)
	})
	http.HandleFunc("/set-default-command", func(w http.ResponseWriter, r *http.Request) {
		handler2.ChangeDefaultCommand(w, r, db)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on %s", config.Port)
	})

	log.Printf("Starting server on %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
