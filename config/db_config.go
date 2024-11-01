package config

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

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
