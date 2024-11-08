package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/templates"
	"github.com/vigneshrajj/gofind/models"

	"gorm.io/gorm"
)

func convertToFileURL(r *http.Request, filePath string) string {
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/files/%s", protocol, r.Host, filepath.Base(filePath))
}

func HandleAddCommand(w http.ResponseWriter, r *http.Request, data []string, db *gorm.DB) {
	if len(data) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. Add command usage:\n#a <alias> <url-with-args> <description(optional)>")
		return
	}

	command := models.Command{
		Alias: data[1],
		Query: data[2],
		Type: models.SearchCommand,
	}
	if len(data) > 3 {
		command.Description = sql.NullString{String: strings.Join(data[3:], " "), Valid: true}
	}

	isFile := strings.HasPrefix(command.Query, "file://")
	if isFile {
		command.Query = convertToFileURL(r, command.Query)
	}

	err := database.CreateCommand(db, command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, err.Error())
		return
	}
	templates.MessageTemplate(w, "Added Command: "+data[1]+", URL: "+data[2])
}

func HandleListCommands(w http.ResponseWriter, data []string, db *gorm.DB) {
	if len(data) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. List command usage: #l")
		return
	}
	commands := database.ListCommands(db)
	templates.ListCommandsTemplate(w, commands)
}

func ChangeDefaultCommand(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	alias := r.URL.Query().Get("default")
	response := "Default Command has been changed successfully to " + alias
	err := database.SetDefaultCommand(db, alias)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = "Error setting default command."
	}
	w.Write([]byte(response))
}

func HandleDeleteCommand(w http.ResponseWriter, data []string, db *gorm.DB) {
	if len(data) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. Delete command usage: #d <alias>")
		return
	}
	command := database.SearchCommand(db, data[1], false)
	if command == (models.Command{}) {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Command not found.")
		return
	}
	if command.Type == models.ApiCommand || command.Type == models.UtilCommand {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Cannot delete built-in utilities or api commands.")
		return
	}
	database.DeleteCommand(db, data[1])
	fmt.Fprintf(w, "Deleted Command: %s", data[1])
}
