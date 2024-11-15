package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/handlers"
	"github.com/vigneshrajj/gofind/internal/templates"
	"gorm.io/gorm"
)

func HandleRoutes(db *gorm.DB) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("/files"))))

	http.HandleFunc("/opensearch.xml", func(w http.ResponseWriter, r *http.Request) {
		templates.OpenSearchDescriptionTemplate(w)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		handlers.HandleQuery(w, r, query, db)
	})


	http.HandleFunc("/opensearch-suggestions", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		handlers.HandleOpenSearchSuggestions(w, query, db)
	})

	http.HandleFunc("/filter_commands", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleFilteredListCommands(w, r, db)
	})

	http.HandleFunc("/set-default-command", func(w http.ResponseWriter, r *http.Request) {
		handlers.ChangeDefaultCommand(w, r, db)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on %s", config.Port)
	})
}

func StartServer(DbPath string, Port string) error {
	_, db, err := database.NewDBConnection(DbPath)
	if err != nil {
		return err
	}

	database.EnsureDefaultCommandsExist(db)
	if config.ItToolsUrl != "" {
		database.EnsureUtilCommandsExist(db)
	}

	if config.EnableAdditionalCommands {
		database.EnsureAdditionalCommandsExist(db)
	}

	HandleRoutes(db)

	log.Printf("Starting server on %s", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
	return nil
}

