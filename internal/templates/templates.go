package templates

import (
	"fmt"
	"github.com/vigneshrajj/gofind/models"
	"html/template"
	"net/http"
)

type ListCommandsPageData struct {
	Commands []models.Command
}

func ListCommandsTemplate(w http.ResponseWriter, commands []models.Command) {
	data := ListCommandsPageData{
		Commands: commands,
	}
	tmpl, err := template.ParseFiles("static/templates/list_commands.html")
	if err != nil {
		fmt.Fprint(w, "ListCommands Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "ListCommands Template couldn't be executed.")
	}
}

type MessagePageData struct {
	Message string
}

func MessageTemplate(w http.ResponseWriter, msg string) {
	data := MessagePageData{
		Message: msg,
	}
	tmpl, err := template.ParseFiles("static/templates/message.html")
	if err != nil {
		fmt.Fprint(w, "MessageTemplate Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "MessageTemplate Template couldn't be executed.")
	}
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
	tmpl, err := template.ParseFiles("static/templates/base64.html")
	if err != nil {
		fmt.Fprint(w, "Base64 Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "Base64 Template couldn't be executed.")
	}
}

func Base64DecodeTemplate(w http.ResponseWriter, decoded string) {
	data := B64PageData{
		Value: decoded,
		Type:  Decoded,
	}
	tmpl, err := template.ParseFiles("static/templates/base64.html")
	if err != nil {
		fmt.Fprint(w, "Base64 Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "Base64 Template couldn't be executed.")
	}
}

type Sha256PageData struct {
	Value string
}

func Sha256Template(w http.ResponseWriter, hashed string) {
	data := Sha256PageData{
		Value: hashed,
	}
	tmpl, err := template.ParseFiles("static/templates/sha256.html")
	if err != nil {
		fmt.Fprint(w, "Sha256 Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "Sha256 Template couldn't be executed.")
	}
}
