package tests

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/handler"
	"github.com/vigneshrajj/gofind/models"
)

func setupCommandHandlerTest() func() {
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

func TestCreateCommand(t *testing.T) {
	defer setupCommandHandlerTest()()
	cmd := models.Command{
		Alias: "help",
		Query: "SELECT * FROM commands",
		Type: models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	if err := db.Create(&cmd).Error; err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	var count int64
	db.Model(&models.Command{}).Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
}

func TestDeleteCommand(t *testing.T) {
	defer setupCommandHandlerTest()()
	cmd := models.Command{
		Alias: "help",
		Query: "SELECT * FROM commands",
		Type: models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	handler.CreateCommand(db, cmd)
	var count int64
	db.Model(&models.Command{}).Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
	handler.DeleteCommand(db, "help")
	db.Model(&models.Command{}).Count(&count)
	if count != 0 {
		t.Fatalf("Expected 0 command, but got %v", count)
	}
}

func TestListCommands(t *testing.T) {
	defer setupCommandHandlerTest()()
	cmd := models.Command{
		Alias: "help",
		Query: "SELECT * FROM commands",
		Type: models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	handler.CreateCommand(db, cmd)
	var count int64
	commands := handler.ListCommands(db)
	count = int64(len(commands))
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
}

func TestPartialSearchCommand(t *testing.T) {
	defer setupCommandHandlerTest()()
	cmd := models.Command{
		Alias:        "goo",
		Query:        "google.com",
		Type:         models.ApiCommand,
		Description:  sql.NullString{String: "Search in google", Valid: true},
	}
	if err := handler.CreateCommand(db, cmd); err != nil {
		t.Fatal(err)
	}
	command, err := handler.SearchCommand(db, "g", true)
	if err != nil {
		t.Fatalf("Failed to search command: %v", err)
	}
	if command == (models.Command{}) {
		t.Fatalf("Expected 1 command, got %d", 0)
	}
}
