package tests

import (
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/server"
	"io"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestServerIsRunning(t *testing.T) {
	go server.StartServer()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:3005")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestDBConnection(t *testing.T) {
	dbSql, _, err := database.NewDBConnection(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := dbSql.Ping(); err != nil {
		t.Fatalf("Database connection is not alive: %v", err)
	}
}
