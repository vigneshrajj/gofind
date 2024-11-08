package models

import (
	"database/sql"
	"time"
)

type CommandType string
const (
	UtilCommand CommandType = "util"
	ApiCommand CommandType = "api"
	SearchCommand CommandType = "search"
)

type Command struct {
	Alias string `gorm:"primaryKey;type:VARCHAR(50);not null" json:"alias"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Query string
	Type CommandType `gorm:"type:VARCHAR(50);not null" json:"type"`
	Description sql.NullString
	IsDefault bool
}
