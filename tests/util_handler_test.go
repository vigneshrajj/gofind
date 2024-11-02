package tests

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/vigneshrajj/gofind/handler"
	"github.com/vigneshrajj/gofind/service"
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

func TestB64EncodeService(t *testing.T) {
	data := "test"
	encoded := service.GetB64(data)
	if encoded != "dGVzdA==" {
		t.Fatalf("Expected dGVzdA==, but got %v", encoded)
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

func TestB64DecodeService(t *testing.T) {
	data := "dGVzdA=="
	decoded := service.GetB64Decode(data)
	if decoded != "test" {
		t.Fatalf("Expected test, but got %v", decoded)
	}
}

func TestSha256UtilCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "sha256 test"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestSha256Service(t *testing.T) {
	data := "test"
	hashed := service.Sha256(data)
	if hashed != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Fatalf("Expected 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08, but got %v", hashed)
	}
}
