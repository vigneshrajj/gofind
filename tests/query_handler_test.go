package tests

import (
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
	config.InsertDefaultCommands(db)
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

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
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

	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestExtraArgsDeleteQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d alias extra"
	w:= httptest.NewRecorder()

	handler.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestLessArgsDeleteQueryCommand(t *testing.T) {
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
	w = httptest.NewRecorder()
	query = "alias"

	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com" {
		t.Fatalf("Expected Location header to be 'https://google.com', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectByPartialMatchQueryCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "al"

	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com" {
		t.Fatalf("Expected Location header to be 'https://google.com', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectByDuplicatePartialMatchCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "#a al https://youtube.com"
	urlEncodedQuery = url.QueryEscape(query)
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	query = "al"

	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://youtube.com" {
		t.Fatalf("Expected Location header to be 'https://youtube.com', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectInvalidQueryWithDefaultCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	w := httptest.NewRecorder()
	query := "invalid"
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+query, nil)

	w = httptest.NewRecorder()
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://www.google.com/search?q=invalid" {
		t.Fatalf("Expected Location header to be 'https://www.google.com/search?q=invalid', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectQueryWithNArgsCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q=%s"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias search some string"

	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=search+some+string" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=search+some+string', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectQueryWithMultipleArgsCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={1}+{2}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias search string"
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=search+string" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=search+string', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectQueryWithOneArgCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={1}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias search"

	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=search" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=search', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectQueryWithKeyValueArgCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={key:val,key2:val2,key3:val3}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias key2"

	handler.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=val2" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=val2', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectQueryWithMultipleKeyValueArgCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={key:val,key2:val2,key3:val3}+{key:val,key2:val2,key3:val3}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias key2 key3"
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=val2+val3" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=val2+val3', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectQueryWithInvalidArgsCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={1}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handler.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias"
	handler.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}
