package tests

import (
	"github.com/vigneshrajj/gofind/internal/handlers"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestB64UtilCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "b64 test"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestD64UtilCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "d64 dGVzdA=="
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestSha256UtilCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "sha256 test"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}
