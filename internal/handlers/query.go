package handlers

import (
	"fmt"
	"github.com/vigneshrajj/gofind/internal/database"
	"github.com/vigneshrajj/gofind/internal/templates"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/gorm"
)

func isKeyValueArg(query string) bool {
	bracketIndex := strings.Index(query, "{")
	if bracketIndex == -1 {
		return false
	}
	colonIndex := strings.Index(query[bracketIndex:], ":")
	if colonIndex == -1 {
		return false
	}
	return true
}

func replaceKeyWithValue(input string, choice string) string {
	re := regexp.MustCompile(`{([^}]*)}`)
	matches := re.FindStringSubmatch(input)

	if len(matches) == 0 {
		return input
	}

	content := matches[1]
	keyValuePairs := strings.Split(content, ",")

	kvMap := make(map[string]string)
	for _, pair := range keyValuePairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			kvMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	if value, found := kvMap[choice]; found {
		return strings.Replace(input, matches[0], value, 1)
	}

	return input
}

func HandleRedirectQuery(w http.ResponseWriter, r *http.Request, data []string, db *gorm.DB) {
	alias := data[0]
	command, err := database.SearchCommand(db, alias, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, err.Error())
		return
	}

	if command.Type == models.UtilCommand {
		HandleUtilCommand(w, data)
		return
	}

	if command == (models.Command{}) {
		var defaultCommand models.Command
		defaultCommand, err = database.GetDefaultCommand(db)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.MessageTemplate(w, "Command not found.")
			return
		}
		command = defaultCommand
		data = append([]string{command.Alias}, data...)
	}

	query := command.Query
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Command not found.")
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

	if isKeyValueArg(query) {
		for i := 1; i < len(data); i++ {
			query = replaceKeyWithValue(query, data[i])
		}
		argCountInQuery := strings.Count(query, "{")
		if argCountInQuery > 0 {
			w.WriteHeader(http.StatusBadRequest)
			templates.MessageTemplate(w, "Invalid arguments provided")
			return
		}
		http.Redirect(w, r, query, http.StatusFound)
		return
	}

	argCount := len(data) - 1
	for i := argCount; i >= 1; i-- {
		query = strings.Replace(query, fmt.Sprintf("{%d}", i), data[i], -1)
		query = strings.Replace(query, fmt.Sprintf("{%d}", i), data[i], -1)
	}

	argCountInQuery := strings.Count(query, "{")
	isNArgQuery := strings.Count(query, "%s") == 1
	if argCountInQuery > 0 && !isNArgQuery {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided")
		return
	}

	http.Redirect(w, r, query, http.StatusFound)
}

func HandleQuery(w http.ResponseWriter, r *http.Request, query string, db *gorm.DB) {
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Query cannot be empty")
		return
	}
	data := strings.Split(query, " ")
	switch data[0] {
	case "#a":
		HandleAddCommand(w, data, db)
	case "#d":
		HandleDeleteCommand(w, data, db)
	case "#l":
		HandleListCommands(w, data, db)
	default:
		HandleRedirectQuery(w, r, data, db)
	}
}
