package templates

import (
	"github.com/vigneshrajj/gofind/models"
	"html/template"
	"net/http"
)

type ListCommandsPageData struct {
	GroupedCommands map[models.CommandType][]models.Command
}

func groupByType(commands []models.Command) map[models.CommandType][]models.Command {
	groupedCommands := make(map[models.CommandType][]models.Command)

	for _, command := range commands {
		groupedCommands[command.Type] = append(groupedCommands[command.Type], command)
	}

	return groupedCommands
}

func ListCommandsTemplate(w http.ResponseWriter, commands []models.Command) {
	data := ListCommandsPageData{
		GroupedCommands: groupByType(commands),
	}
	tmpl := template.Must(template.ParseFiles("static/templates/list_commands.html"))
	tmpl.Execute(w, data)
}

type MessagePageData struct {
	Message string
}

func MessageTemplate(w http.ResponseWriter, msg string) {
	data := MessagePageData{
		Message: msg,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/message.html"))
	tmpl.Execute(w, data)
}

type B64PageType string

const (
	Encoded B64PageType = "encoded"
	Decoded B64PageType = "decoded"
)

type B64PageData struct {
	Value string
	Type  B64PageType
}

func Base64Template(w http.ResponseWriter, encoded string) {
	data := B64PageData{
		Value: encoded,
		Type:  Encoded,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/base64.html"))
	tmpl.Execute(w, data)
}

func Base64DecodeTemplate(w http.ResponseWriter, decoded string) {
	data := B64PageData{
		Value: decoded,
		Type:  Decoded,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/base64.html"))
	tmpl.Execute(w, data)
}

type Sha256PageData struct {
	Value string
}

func Sha256Template(w http.ResponseWriter, hashed string) {
	data := Sha256PageData{
		Value: hashed,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/sha256.html"))
	tmpl.Execute(w, data)
}
