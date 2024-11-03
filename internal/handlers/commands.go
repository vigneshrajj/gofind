package handlers

import (
	"database/sql"
	"fmt"
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/templates"
	"github.com/vigneshrajj/gofind/models"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func HandleAddCommand(w http.ResponseWriter, data []string, db *gorm.DB) {
	if len(data) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. Add command usage:\n#a <alias> <url-with-args> <description(optional)>")
		return
	}

	command := models.Command{
		Alias: data[1],
		Query: data[2],
	}
	if len(data) > 3 {
		command.Description = sql.NullString{String: strings.Join(data[3:], " "), Valid: true}
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
	command, err := database.SearchCommand(db, data[1], false)
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
	err = database.DeleteCommand(db, data[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, err.Error())
		return
	}
	fmt.Fprintf(w, "Deleted Command: %s", data[1])
}
