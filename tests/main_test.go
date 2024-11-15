package tests

import (
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/server"

	"gorm.io/gorm"
)

var db *gorm.DB
var (
	wg          sync.WaitGroup
	serverReady = false
	once        sync.Once
)

func setupServerTest() func() {
	os.Symlink("../static", "./static")
	once.Do(func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.StartServer(":memory:", ":3005")
		}()
		// Wait for the server to start
		time.Sleep(100 * time.Millisecond)
		serverReady = true
	})

	for !serverReady {
		time.Sleep(10 * time.Millisecond)
	}

	return func() {}
}

func TestServerIsRunning(t *testing.T) {
	defer setupServerTest()()

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


func TestServerFailWithInvalidDbPath(t *testing.T) {
	err := server.StartServer("db/db/invalid.db", ":3005")

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
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

func TestDBConnectionWithInvalidFileName(t *testing.T) {
	_, _, err := database.NewDBConnection("db/db/invalid.db")
	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
}

func TestSearchEndpoint(t *testing.T) {
	defer setupServerTest()()

	resp, err := http.Get("http://localhost:3005/search?query=g")
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

func TestSetDefaultCommandEndpoint(t *testing.T) {
	defer setupServerTest()()
	resp, err := http.Get("http://localhost:3005/set-default-command?default=g")
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
