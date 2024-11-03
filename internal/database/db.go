package database

import (
	"database/sql"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDBConnection(dbFileName string) (*sql.DB, *gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	EnsureCommandTableExists(db)

	dbSql, _ := db.DB()
	return dbSql, db, nil
}

func EnsureCommandTableExists(db *gorm.DB) {
	db.AutoMigrate(&models.Command{})
}

func EnsureDefaultCommandsExist(db *gorm.DB) {
	defaultCommands := []models.Command{
		{
			Alias:       "g",
			Query:       "https://www.google.com/search?q=%s",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "Google Search", Valid: true},
			IsDefault:   true,
		},
		{
			Alias:       "b64",
			Query:       "",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "Convert To Base 64 string", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:       "d64",
			Query:       "",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "Decode Base 64 string", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:       "sha256",
			Query:       "",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "Convert To SHA 256 string", Valid: true},
			IsDefault:   false,
		},
	}
	for _, command := range defaultCommands {
		FirstOrCreateCommand(db, command)
	}
}

func EnsureAdditionalCommandsExist(db *gorm.DB) {
	additionalCommands := []models.Command{
		{
			Alias:       "y",
			Query:       "https://www.youtube.com/results?search_query=%s",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "Youtube", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:       "ddg",
			Query:       "https://duckduckgo.com/?q=%s",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "DuckDuckGo", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "ddl",
			Query:       "https://lite.duckduckgo.com/lite/?q=%s",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "DuckDuckGo Lite", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "gh",
			Query:       "https://github.com/search?q=%s&type=repositories",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "Github Repos", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "npm",
			Query:       "https://www.npmjs.com/search?q=%s",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "Node Package Manager (NPM)", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "m",
			Query:       "https://mail.google.com/mail/u/{r:0,vr:1}/#inbox",
			Type:        models.SearchCommand,
			Description: sql.NullString{String: "GMail", Valid: true},
			IsDefault:   false,
		},
	}
	for _, command := range additionalCommands {
		FirstOrCreateCommand(db, command)
	}
}
