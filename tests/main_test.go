package tests

import (
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestServerIsRunning(t *testing.T) {
	go config.StartServer()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:3005")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestDBConnection(t *testing.T) {
	dbSql, _, err := config.NewDBConnection(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := dbSql.Ping(); err != nil {
		t.Fatalf("Database connection is not alive: %v", err)
	}
}
