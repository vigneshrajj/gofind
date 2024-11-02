package handler

import (
	"net/http"

	"gorm.io/gorm"
)

func ChangeDefaultCommand(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	alias := r.URL.Query().Get("default")
	response := "Default Command has been changed successfully to "+ alias
	SetDefaultCommand(db, alias)
	w.Write([]byte(response))
}
