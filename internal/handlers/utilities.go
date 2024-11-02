package handlers

import (
	"net/http"

	"github.com/vigneshrajj/gofind/internal/helpers"
	"github.com/vigneshrajj/gofind/internal/templates"
)

func HandleUtilCommand(w http.ResponseWriter, data []string) {
	alias := data[0]
	switch alias {
	case "b64":
		HandleB64Util(w, data)
	case "d64":
		HandleD64Util(w, data)
	case "sha256":
		HandleSha256Util(w, data)
	default:
		HandleUtilNotFound(w)
	}
}

func HandleUtilNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	templates.MessageTemplate(w, "Command not found.")
}

func HandleSha256Util(w http.ResponseWriter, data []string) {
	if len(data) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. SHA256 encode command usage: sha256 <string>")
		return
	}
	encoded := helpers.Sha256(data[1])
	templates.Sha256Template(w, encoded)
}

func HandleD64Util(w http.ResponseWriter, data []string) {
	if len(data) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. Base64 Decode command usage: d64 <string>")
		return
	}
	decoded := helpers.GetB64Decode(data[1])
	templates.Base64DecodeTemplate(w, decoded)
}

func HandleB64Util(w http.ResponseWriter, data []string) {
	if len(data) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		templates.MessageTemplate(w, "Invalid number of arguments provided. Base64 command usage: b64 <string>")
		return
	}
	encoded := helpers.GetB64(data[1])
	templates.Base64Template(w, encoded)
}
