package server

import (
	"fmt"
	"github.com/vigneshrajj/gofind/internal/database"
	handler2 "github.com/vigneshrajj/gofind/internal/handlers"
	"log"
	"net/http"
)

const DbName = "db/gofind.db"

func StartServer() {
	_, db, err := database.NewDBConnection(DbName)
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
		fmt.Fprint(w, "Server running on :3005")
	})

	log.Print("Starting server on :3005")
	log.Fatal(http.ListenAndServe(":3005", nil))
}
