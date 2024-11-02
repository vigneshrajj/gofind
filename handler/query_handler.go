package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/gorm"
)

func DeleteQuery(w http.ResponseWriter, data []string, db *gorm.DB) {
		if len(data) != 2 {
			http.Error(w, "Invalid number of arguments provided. Delete command usage:\n#d <alias>", http.StatusBadRequest)
			return
		}
		command, err := SearchCommand(db, data[1], false)
		if command == (models.Command{}) {
			http.Error(w, "Command not found.", http.StatusBadRequest)
			return
		}
		err = DeleteCommand(db, data[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Deleted Command: %s", data[1])
}

func ListQuery(w http.ResponseWriter, data []string, db *gorm.DB) {
		if len(data) != 1 {
			http.Error(w, "Invalid number of arguments provided. List command usage:\n#l", http.StatusBadRequest)
			return
		}
		commands := ListCommands(db)
		commandsJson, err := json.Marshal(commands)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Commands: %s", string(commandsJson))
}

func AddQuery(w http.ResponseWriter, data []string, db *gorm.DB) {
		if len(data) < 3 {
			http.Error(w, "Invalid number of arguments provided. Add command usage:\n#a <alias> <url-with-args> <description(optional)>", http.StatusBadRequest)
			return
		}

		command := models.Command{
			Alias: data[1],
			Query: data[2],
		}
		if len(data) > 3 {
			command.Description = sql.NullString{String: strings.Join(data[3:], " "), Valid: true}
		}
			
		err := CreateCommand(db, command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Added Command: %s\nURL: %s", data[1], data[2])
}

func RedirectQuery(w http.ResponseWriter, r *http.Request, data []string, db *gorm.DB) {
	alias := data[0]
	command, err := SearchCommand(db, alias, true)
	if err != nil {
		http.Error(w, "Command not found", http.StatusBadRequest)
		return
	}
	query := command.Query
	if query == "" {
		http.Error(w, "Command not found", http.StatusBadRequest)
		return
	}

	startsWithHttp := strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://")
	if !startsWithHttp {
		query = "https://" + query
	}

	if strings.Contains(query, "%s") {
		query = fmt.Sprintf(query, url.QueryEscape(strings.Join(data[1:], " ")))
		http.Redirect(w, r, query, http.StatusFound)
		return
	}

	argCount := len(data) - 1
	for i := argCount; i >= 1; i-- {
		query = strings.Replace(query, fmt.Sprintf("$%d", i), data[i], -1)
	}

	argCountInQuery := strings.Count(query, "$")
	isNArgQuery := strings.Count(query, "%s") == 1
	if argCountInQuery > 0 && !isNArgQuery {
		http.Error(w, "Invalid number of arguments provided", http.StatusBadRequest)
	}

	http.Redirect(w, r, query, http.StatusFound)
}

func HandleQuery(w http.ResponseWriter, r *http.Request, query string, db *gorm.DB) {
	if query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}
	data := strings.Split(query, " ")
	switch data[0] {
	case "#a":
		AddQuery(w, data, db)
	case "#d":
		DeleteQuery(w, data, db)
	case "#l":
		ListQuery(w, data, db)
	default:
		RedirectQuery(w, r, data, db)
	}
}
