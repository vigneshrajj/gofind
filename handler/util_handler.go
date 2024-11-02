package handler

import (
	"net/http"

	"github.com/vigneshrajj/gofind/service"
	"gorm.io/gorm"
)

func HandleUtilCommand(w http.ResponseWriter, r *http.Request, data []string, db *gorm.DB) {
	alias := data[0]
	switch alias {
	case "b64":
		if len(data) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			service.MessagePage(w, "Invalid number of arguments provided. Base64 command usage: b64 <string>")
			return
		}
		encoded := service.GetB64(data[1])
		service.Base64Page(w, encoded)
	case "d64":
		if len(data) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			service.MessagePage(w, "Invalid number of arguments provided. Base64 Decode command usage: d64 <string>")
			return
		}
		decoded := service.GetB64Decode(data[1])
		service.Base64DecodePage(w, decoded)
	case "sha256":
		if len(data) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			service.MessagePage(w, "Invalid number of arguments provided. SHA256 encode command usage: sha256 <string>")
			return
		}
		encoded := service.Sha256(data[1])
		service.Sha256Page(w, encoded)
	default:
		w.WriteHeader(http.StatusBadRequest)
		service.MessagePage(w, "Command not found.")
	}
}
