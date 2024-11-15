package templates

import (
	"html/template"
	text_template "text/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/models"
)

type ArgType string

const (
	KeyVal ArgType = "keyval"
	Num ArgType = "num"
	Any ArgType = "any"
	None ArgType = "none"
)

type CommandWithArgs struct {
	models.Command

	QueryHostname string
	ArgType ArgType
	ArgsKeyVal map[string]string
	ArgsNum []int
}

type ListCommandsPageData struct {
	HostUrl string
	Type 	 models.CommandType
	EnableUtils bool
}

func ExtractNumArgs(query string) []int {
	re := regexp.MustCompile(`\{(\d+)\}`)
	matches := re.FindAllStringSubmatch(query, -1)

	var result []int

	for _, match := range matches {
		num, err := strconv.Atoi(match[1])
		if err == nil {
			result = append(result, num)
		} else {
			log.Fatal(err.Error())
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
				key := strings.Split(kv[0], "$(default)")[0]
				result[key] = kv[1]
			}
		}
	}

	return result
}


func extractHostnameFromQuery(query string) string {
	startsWithValidProtocol := strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://")

	if !startsWithValidProtocol {
		query = "https://" + query
	}
	// regexp for getting string between second / and third / or end of line
	re := regexp.MustCompile(`\/([^\/]+)\/?`)
	matches := re.FindAllStringSubmatch(query, -1)

	if len(matches) < 1 || len(matches[0]) < 2 {
		return ""
	}

	return matches[0][1]
}

	
func addArgsAndHostToCommands(commands []models.Command) []CommandWithArgs {
	newCommands := make([]CommandWithArgs, 0, len(commands))
	for _, command := range commands {
		uri, err := url.Parse(command.Query)
		if err != nil {
			newCommands = append(newCommands, CommandWithArgs{
				Command: command,
				QueryHostname: command.Query,
				ArgType: None,
			})
			continue
		}
		commandWithArgs := CommandWithArgs{
			Command: command,
			ArgsNum: ExtractNumArgs(command.Query),
			ArgsKeyVal: ExtractKeyValArgs(command.Query),
		}

		isAnyArg := strings.Count(commandWithArgs.Query, "%s") == 1
		if len(commandWithArgs.ArgsKeyVal) > 0 {
			commandWithArgs.ArgType = KeyVal
		} else if (len(commandWithArgs.ArgsNum) > 0) {
			commandWithArgs.ArgType = Num
		} else if isAnyArg {
			commandWithArgs.ArgType = Any
		} else {
			commandWithArgs.ArgType = None
		}

		hostname := uri.Hostname()
		commandWithArgs.QueryHostname = hostname

		if hostname == "" {
			commandWithArgs.QueryHostname = extractHostnameFromQuery(command.Query)
		}
		newCommands = append(newCommands, commandWithArgs)
	}

	return newCommands
}


var helpers template.FuncMap = map[string]interface{}{
	"isLast": func(index int, len int) bool {
		return index+1 == len
	},
}

func ListCommandsTemplate(w http.ResponseWriter, command_type models.CommandType) {
	data := ListCommandsPageData{
		HostUrl: config.HostUrl,
		Type: command_type,
		EnableUtils: config.ItToolsUrl != "",
	}
	tmpl := template.Must(template.New("list_commands.html").Funcs(helpers).ParseFiles("static/templates/list_commands.html", "static/templates/filtered_commands_list.html", "static/templates/command_tabs.html"))
	tmpl.Execute(w, data)
}

type FilteredCommandsPageData struct {
	Offset int
	Commands []CommandWithArgs
}

func FilteredListCommandsTemplate(w http.ResponseWriter, commands []models.Command, offset int) {
	data := FilteredCommandsPageData{
		Offset: offset + len(commands),
		Commands: addArgsAndHostToCommands(commands),
	}

	tmpl := template.Must(template.New("filtered_commands_list.html").Funcs(helpers).ParseFiles("static/templates/filtered_commands_list.html"))
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

type MultiQueryPageData struct {
	Queries []string
}

func MultiQueryTemplate(w http.ResponseWriter, queries []string) {
	data := MultiQueryPageData{
		Queries: queries,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/multi_query.html"))
	tmpl.Execute(w, data)
}

type NotificationData struct {
	Title string
}

func NotificationTemplate(w http.ResponseWriter, title string) {
	data := NotificationData{
		Title: title,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/notification.html"))
	tmpl.Execute(w, data)
}

type OpenSearchPageData struct {
	HostUrl string
}

func OpenSearchDescriptionTemplate(w http.ResponseWriter) {
	data := OpenSearchPageData{
		config.HostUrl,
	}
	tmpl := text_template.Must(text_template.ParseFiles("static/opensearch.xml"))
	tmpl.Execute(w, data)
}
