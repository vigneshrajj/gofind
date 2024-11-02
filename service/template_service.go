package service

import (
	"fmt"
	"html/template"
	"net/http"
)

type MainLayoutData struct {
	PageTitle string
}

func MainLayout(w http.ResponseWriter) {
	data := MainLayoutData{
		PageTitle: "Test title",
	}
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		fmt.Fprint(w, "MainLayout Template not found.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "MainLayout Template couldn't be executed.")
	}
}
