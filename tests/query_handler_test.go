package tests

import (
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/handler"
	"github.com/vigneshrajj/gofind/models"
)

func setupQueryHandlerTest() func() {
	var err error
	_, db, err = config.NewDBConnection(":memory:")
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Command{}); err != nil {
		panic(err)
	}

	return func() {}
}

func TestEmptyQuery(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := ""
	w:= httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
	if string(body) != "Query cannot be empty\n" {
		t.Fatalf("Expected body to be 'Query cannot be empty', but got %v", string(body))
	}
}

func TestAddQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://url"
	w:= httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestDuplicateAddQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://url"
	w := httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	w = httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestInvalidAddQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias"
	w:= httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestDeleteNonExistingQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d alias"
	w:= httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestDeleteQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	w := httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)

	query = "#d alias"
	w = httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestExtraParamsDeleteQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d alias extra"
	w:= httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestLessParamsDeleteQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d"
	w:= httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestListQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#l"
	w := httptest.NewRecorder()
	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestRedirectQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	query = "alias"
	w = httptest.NewRecorder()
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com" {
		t.Fatalf("Expected Location header to be 'https://google.com', but got %v", resp.Header.Get("Location"))
	}
}
