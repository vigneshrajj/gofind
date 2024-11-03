package tests

import (
	"database/sql"
	"github.com/vigneshrajj/gofind/internal/database"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/models"
)

func setupCommandsOperationsTest() func() {
	var err error
	_, db, err = database.NewDBConnection(":memory:")
	if err != nil {
		panic(err)
	}

	return func() {}
}

func TestCreateCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "some query",
		Type:        models.UtilCommand,
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

func TestFirstOrCreateCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "some query",
		Type:        models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	if err := database.FirstOrCreateCommand(db, cmd); err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	var count int64
	db.Model(&models.Command{}).Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
}

func TestFirstOrCreateCommandWithExistingCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "some query",
		Type:        models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	if err := database.FirstOrCreateCommand(db, cmd); err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	if err := database.FirstOrCreateCommand(db, cmd); err != nil {
		t.Fatalf("Failed to create command again: %v", err)
	}
	var count int64
	db.Model(&models.Command{}).Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
}
	

func TestDeleteCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "https://google.com",
		Type:        models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	err := database.CreateCommand(db, cmd)
	if err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	var count int64
	db.Model(&models.Command{}).Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
	err = database.DeleteCommand(db, "help")
	if err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	db.Model(&models.Command{}).Count(&count)
	if count != 0 {
		t.Fatalf("Expected 0 command, but got %v", count)
	}
}

func TestListCommands(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "https://google.com",
		Type:        models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	err := database.CreateCommand(db, cmd)
	if err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	var count int64
	commands := database.ListCommands(db)
	count = int64(len(commands))
	if count != 1 {
		t.Fatalf("Expected 1 command, but got %v", count)
	}
}

func TestPartialSearchCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "goo",
		Query:       "google.com",
		Type:        models.ApiCommand,
		Description: sql.NullString{String: "Search in google", Valid: true},
	}
	if err := database.CreateCommand(db, cmd); err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	command, err := database.SearchCommand(db, "g", true)
	if err != nil {
		t.Fatalf("Failed to search command: %v", err)
	}
	if command == (models.Command{}) {
		t.Fatalf("Expected 1 command, got %d", 0)
	}
}

func TestGetDefaultCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "https://google.com",
		Type:        models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
		IsDefault:   true,
	}
	err := database.CreateCommand(db, cmd)
	if err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	command, err := database.GetDefaultCommand(db)
	if err != nil {
		t.Fatalf("Failed to get default command: %v", err)
	}
	if command == (models.Command{}) {
		t.Fatalf("Expected 1 command, got %d", 0)
	}
}

func TestGetDefaultCommand_Error(t *testing.T) {
	defer setupCommandsOperationsTest()()
	command, _ := database.GetDefaultCommand(db)
	if command != (models.Command{}) {
		t.Fatalf("Expected 0 command, got %d", 1)
	}
}

func TestSetDefaultCommand(t *testing.T) {
	defer setupCommandsOperationsTest()()
	cmd := models.Command{
		Alias:       "help",
		Query:       "https://google.com",
		Type:        models.UtilCommand,
		Description: sql.NullString{String: "List all available commands", Valid: true},
	}
	err := database.CreateCommand(db, cmd)
	if err != nil {
		t.Fatalf("Failed to create a command: %v", err)
	}
	err = database.SetDefaultCommand(db, "help")
	if err != nil {
		t.Fatalf("Failed to set default command: %v", err)
	}
	command, err := database.GetDefaultCommand(db)
	if err != nil {
		t.Fatalf("Failed to get default command: %v", err)
	}
	if command == (models.Command{}) {
		t.Fatalf("Expected 1 command, got %d", 0)
	}
}
