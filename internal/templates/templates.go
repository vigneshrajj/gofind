package templates

import (
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/vigneshrajj/gofind/models"
)

type ArgType string

const (
	KeyVal ArgType = "keyval"
	Num ArgType = "num"
	Any ArgType = "any"
)

type CommandWithArgs struct {
	models.Command

	QueryHostname string
	ArgType ArgType
	ArgsKeyVal map[string]string
	ArgsNum []int
}

type ListCommandsPageData struct {
	GroupedCommands map[models.CommandType][]CommandWithArgs
}

func ExtractNumArgs(query string) []int {
	re := regexp.MustCompile(`\{(\d+)\}`)
	matches := re.FindAllStringSubmatch(query, -1)

	var result []int

	for _, match := range matches {
		num, err := strconv.Atoi(match[1])
		if err == nil {
			result = append(result, num)
		}
	}

	return result
}

func ExtractKeyValArgs(query string) map[string]string {
	re := regexp.MustCompile(`\{([^\}]+)\}`)

	matches := re.FindAllStringSubmatch(query, -1)

	result := make(map[string]string)

	for _, match := range matches {
		keyValuePairs := match[1]
		pairs := strings.Split(keyValuePairs, ",")
		for _, pair := range pairs {
			kv := strings.Split(pair, ":")
			if len(kv) == 2 {
				result[kv[0]] = kv[1]
			}
		}
	}

	return result
}

		

func groupByType(commands []models.Command) map[models.CommandType][]CommandWithArgs {
	groupedCommands := make(map[models.CommandType][]CommandWithArgs)

	for _, command := range commands {
		uri, err := url.Parse(command.Query)
		if err != nil {
			return groupedCommands
		}
		commandWithArgs := CommandWithArgs{
			Command: command,
			ArgsNum: ExtractNumArgs(command.Query),
			ArgsKeyVal: ExtractKeyValArgs(command.Query),
		}
		if len(commandWithArgs.ArgsKeyVal) > 0 {
			commandWithArgs.ArgType = KeyVal
		} else if (len(commandWithArgs.ArgsNum) > 0) {
			commandWithArgs.ArgType = Num
		} else {
			commandWithArgs.ArgType = Any
		}

		commandWithArgs.QueryHostname = uri.Hostname()
		groupedCommands[command.Type] = append(groupedCommands[command.Type], commandWithArgs)
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
