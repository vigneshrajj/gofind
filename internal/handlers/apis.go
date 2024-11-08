package handlers

import (
	"net/http"

	"github.com/vigneshrajj/gofind/internal/helpers"
	"github.com/vigneshrajj/gofind/internal/templates"
)

func HandleApiCommands(w http.ResponseWriter, data []string) {
	alias := data[0]
	switch alias {
	case "todo":
		HandleTodoistApi(w, data)
	}
}

func HandleTodoistApi(w http.ResponseWriter, data []string) {
	if len(data) <= 2 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. Todoist command usage: todo add <task-description> pri:<none|low|med|high> due:<due-date-in-natural-language> labels:[cat,dog]")
		return
	}
	encoded := helpers.Sha256(data[1])
	templates.Sha256Template(w, encoded)
}
