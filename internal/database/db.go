package database

import (
	"database/sql"

	"github.com/vigneshrajj/gofind/config"
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

func EnsureUtilCommandsExist(db *gorm.DB) {
	utilCommands := []models.Command{
		{
			Alias:       "!it",
			Query:       config.ItToolsUrl,
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:       "!bcrypt",
			Query:       config.ItToolsUrl+"/bcrypt?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Generate Bcrypt", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:       "!b64",
			Query:       config.ItToolsUrl+"/base64-string-converter?defaultText={2}&shouldEncode={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Base 64 Encoded Decoder", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!case",
			Query:       config.ItToolsUrl+"/case-converter?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Case Converter", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!color",
			Query:       config.ItToolsUrl+"/color-converter?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Color Converter", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!dt",
			Query:       config.ItToolsUrl+"/date-converter?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Date Converter", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!emoji",
			Query:       config.ItToolsUrl+"/emoji-picker?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Search Emoji", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!hash",
			Query:       config.ItToolsUrl+"/hash-text?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Hash Text", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!http",
			Query:       config.ItToolsUrl+"/http-status-codes?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: HTTP Status Code Lookup", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!jwt",
			Query:       config.ItToolsUrl+"/jwt-parser?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: JWT Parser", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!qr",
			Query:       config.ItToolsUrl+"/qr-code-generator?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: QR Code Generator", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!slug",
			Query:       config.ItToolsUrl+"/slugify-string?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Slugify String", Valid: true},
			IsDefault:   false,
		},
		{
			Alias:			 "!url",
			Query:       config.ItToolsUrl+"/url-parser?defaultText={1}",
			Type:        models.UtilCommand,
			Description: sql.NullString{String: "IT Tools: Parse URL", Valid: true},
			IsDefault:   false,
		},
	}
	for _, command := range utilCommands {
		FirstOrCreateCommand(db, command)
	}
}
