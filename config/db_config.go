package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartServer() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Missing 'query' parameter", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Query: %s", query)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Server running on :3005")
	})

	log.Print("Starting server on :3005")
	log.Fatal(http.ListenAndServe(":3005", nil))
}

func NewDBConnection(dbFileName string) (*sql.DB, *gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	dbSql, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	HandleCommandTable(db)

	return dbSql, db, nil
}

func HandleCommandTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Command{}); err != nil {
		log.Fatalf("Failed to migrate the Command schema: %v", err)
		return err
	}
	return nil
}
