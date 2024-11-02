package tests

import (
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/handlers"
	"net/http/httptest"
	"net/url"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/models"
)

func setupQueryHandlerTest() func() {
	var err error
	_, db, err = database.NewDBConnection(":memory:")
	if err != nil {
		panic(err)
	}
	err = database.EnsureDefaultCommandsExist(db)
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(&models.Command{}); err != nil {
		panic(err)
	}

	return func() {}
}

func TestEmptyQuery(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := ""
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestAddCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://url"
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestAddDuplicateCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://url"
	w := httptest.NewRecorder()
	handlers.HandleQuery(w, nil, query, db)
	w = httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestInvalidAddCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias"
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestDeleteNonExistingCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d alias"
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestDeleteCommandQuery(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	w := httptest.NewRecorder()
	handlers.HandleQuery(w, nil, query, db)
	query = "#d alias"

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestExtraArgsDeleteCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d alias extra"
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestLessArgsDeleteCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#d"
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestListCommandsQuery(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#l"
	w := httptest.NewRecorder()

	handlers.HandleQuery(w, nil, query, db)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
}

func TestRedirectQuery(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com" {
		t.Fatalf("Expected Location header to be 'https://google.com', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectByPartialMatch(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "al"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com" {
		t.Fatalf("Expected Location header to be 'https://google.com', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectBySupercedingPartialMatch(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "#a al https://youtube.com"
	urlEncodedQuery = url.QueryEscape(query)
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	query = "al"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://youtube.com" {
		t.Fatalf("Expected Location header to be 'https://youtube.com', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectNonExistingAliasToDefaultCommand(t *testing.T) {
	defer setupQueryHandlerTest()()
	w := httptest.NewRecorder()
	query := "invalid"
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+query, nil)

	w = httptest.NewRecorder()
	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://www.google.com/search?q=invalid" {
		t.Fatalf("Expected Location header to be 'https://www.google.com/search?q=invalid', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectWithNArgs(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q=%s"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias search some string"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=search+some+string" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=search+some+string', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectWithMultipleArgs(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={1}+{2}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias search string"
	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=search+string" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=search+string', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectWithOneArg(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={1}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias search"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=search" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=search', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectWithKeyValueArg(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={key:val,key2:val2,key3:val3}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias key2"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=val2" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=val2', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectWithMultipleKeyValueArg(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={key:val,key2:val2,key3:val3}+{key:val,key2:val2,key3:val3}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias key2 key3"
	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()
	if resp.StatusCode != 302 {
		t.Fatalf("Expected status code 302, but got %v", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "https://google.com/search?q=val2+val3" {
		t.Fatalf("Expected Location header to be 'https://google.com/search?q=val2+val3', but got %v", resp.Header.Get("Location"))
	}
}

func TestRedirectWithInvalidKeyValueArg(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={key:val,key2:val2,key3:val3}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias abc"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}

func TestRedirectWithInvalidArgs(t *testing.T) {
	defer setupQueryHandlerTest()()
	query := "#a alias https://google.com/search?q={1}"
	urlEncodedQuery := url.QueryEscape(query)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:3005/search?query="+urlEncodedQuery, nil)
	handlers.HandleQuery(w, r, query, db)
	w = httptest.NewRecorder()
	query = "alias"

	handlers.HandleQuery(w, r, query, db)
	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got %v", resp.StatusCode)
	}
}
