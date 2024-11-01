package tests

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/models"
	"github.com/vigneshrajj/gofind/handler"
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
	query := ""
	res, err := handler.HandleQuery(query)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	if res != "" {
		t.Fatalf("Expected empty string, got %s", res)
	}
}
