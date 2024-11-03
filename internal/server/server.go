package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/handlers"
	"gorm.io/gorm"
)

func HandleRoutes(db *gorm.DB) {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		handlers.HandleQuery(w, r, query, db)
	})
	http.HandleFunc("/set-default-command", func(w http.ResponseWriter, r *http.Request) {
		handlers.ChangeDefaultCommand(w, r, db)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on %s", config.Port)
	})
}

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

	if config.EnableAdditionalCommands {
		err = database.EnsureAdditionalCommandsExist(db)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	HandleRoutes(db)

	log.Printf("Starting server on %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
