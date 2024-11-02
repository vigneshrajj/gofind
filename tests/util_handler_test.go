package tests

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/vigneshrajj/gofind/handler"
)

func TestB64UtilCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "b64 test"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)

	handler.HandleQuery(w, r, query, db)
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
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}
