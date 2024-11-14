package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/templates"
	"github.com/vigneshrajj/gofind/models"

	"gorm.io/gorm"
)

func GetHostFromRequest(r *http.Request) string {
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s", protocol, r.Host)
}

func convertToFileURL(r *http.Request, filePath string) string {
	return fmt.Sprintf("%s/files/%s", GetHostFromRequest(r), filepath.Base(filePath))
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
	templates.ListCommandsTemplate(w, models.SearchCommand)
}

func HandleFilteredListCommands(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	page_size := 10
	offset := 0
	searchQuery := r.URL.Query().Get("search_query")
	command_type := r.URL.Query().Get("command_type")

	if r.URL.Query().Get("page_size") != "" {
		var err error
		page_size, err = strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.MessageTemplate(w, "Invalid page_size provided.")
			return
		}
	}
	if r.URL.Query().Get("offset") != "" {
		var err error
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.MessageTemplate(w, "Invalid offset provided.")
			return
		}
	}

	commands, err := database.FilteredListCommands(db, searchQuery, page_size, offset, command_type)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Could not fetch commands." + err.Error())
		return
	}

	templates.FilteredListCommandsTemplate(w, *commands, offset)
}

func ChangeDefaultCommand(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	alias := r.URL.Query().Get("default")
	response := "Default Command has been changed successfully to " + alias
	err := database.SetDefaultCommand(db, alias)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = "Error setting default command."
	}
	templates.NotificationTemplate(w, response)
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

func HandleUserCommands(w http.ResponseWriter, data []string, db *gorm.DB) {
	if len(data) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided")
		return
	}
	name := data[1]
	argument := strings.Join(data[2:], " ")
	scriptsDir := config.ScriptsPath

	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Couldn't read scripts directory: " + err.Error())
		return
	}
	var scriptPath string

	for _, file := range files {
		baseName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		if baseName == name {
			scriptPath = filepath.Join(scriptsDir, file.Name())
			break
		}
	}
	if scriptPath == "" {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Couldn't find the specified script.")
		return
	}

	cmd := exec.Command(scriptPath, argument)
	// cmd.Env = append(os.Environ(), "SHELL=/bin/bash")
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Couldn't execute the script: " + err.Error())
		return
	}

	templates.MessageTemplate(w, "Executed the script: " + string(output))
}
