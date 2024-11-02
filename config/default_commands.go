package config

import (
	"database/sql"
	"github.com/vigneshrajj/gofind/handler"
	"github.com/vigneshrajj/gofind/models"
	"gorm.io/gorm"
)

func InsertDefaultCommands(db *gorm.DB) error {
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
		if err := handler.FirstOrCreateCommand(db, command); err != nil {
			anyerr = err
		}
	}
	if anyerr != nil {
		return anyerr
	}
	return nil
}
