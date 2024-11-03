package database

import (
	"database/sql"
	"log"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDBConnection(dbFileName string) (*sql.DB, *gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	dbSql, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	err = EnsureCommandTableExists(db)
	if err != nil {
		return nil, nil, err
	}

	return dbSql, db, nil
}

func EnsureCommandTableExists(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Command{}); err != nil {
		log.Fatalf("Failed to migrate the Command schema: %v", err)
		return err
	}
	return nil
}

func EnsureDefaultCommandsExist(db *gorm.DB) error {
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
	var anyerr error
	for _, command := range defaultCommands {
		if err := FirstOrCreateCommand(db, command); err != nil {
			anyerr = err
		}
	}
	if anyerr != nil {
		return anyerr
	}
	return nil
}

func EnsureAdditionalCommandsExist(db *gorm.DB) error {
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
	var anyerr error
	for _, command := range additionalCommands {
		if err := FirstOrCreateCommand(db, command); err != nil {
			anyerr = err
		}
	}
	if anyerr != nil {
		return anyerr
	}
	return nil
}
