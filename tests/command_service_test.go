package tests

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/models"
	"github.com/vigneshrajj/gofind/handler"
)

func setupCommandServiceTest() func() {
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

func TestSearchCommand(t *testing.T) {
	defer setupCommandServiceTest()()
	cmd := models.Command{
		Alias:        "g",
		Query:        "google.com",
		Type:         models.ApiCommand,
		Description:  sql.NullString{String: "Search in google", Valid: true},
	}
	if err := handler.CreateCommand(db, cmd); err != nil {
		t.Fatal(err)
	}
	commands := handler.ListCommands(db)
	if len(commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(commands))
	}
}
