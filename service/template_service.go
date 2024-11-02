package service

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/vigneshrajj/gofind/models"
)

type MainLayoutData struct {
	PageTitle string
}

func MainLayout(w http.ResponseWriter) {
	data := MainLayoutData{
		PageTitle: "Test title",
	}
	tmpl, err := template.ParseFiles("static/templates/layout.html")
	if err != nil {
		fmt.Fprint(w, "MainLayout Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "MainLayout Template couldn't be executed.")
	}
}

type ListCommandsPageData struct {
	Commands []models.Command
}

func ListCommandsPage(w http.ResponseWriter, commands []models.Command) {
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

func MessagePage(w http.ResponseWriter, msg string) {
	data := MessagePageData{
		Message: msg,
	}
	tmpl, err := template.ParseFiles("static/templates/message.html")
	if err != nil {
		fmt.Fprint(w, "MessagePage Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "MessagePage Template couldn't be executed.")
	}
}

type B64PageType string
const (
	Encoded B64PageType = "encoded"
	Decoded B64PageType = "decoded"
)

type B64PageData struct {
	Value string
	Type B64PageType
}


func Base64Page(w http.ResponseWriter, encoded string) {
	data := B64PageData{
		Value: encoded,
		Type: "encoded",
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

func Base64DecodePage(w http.ResponseWriter, decoded string) {
	data := B64PageData{
		Value: decoded,
		Type: "decoded",
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
