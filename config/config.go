package config

import "os"


var DbPath = "db/gofind.db"
var Port = ":3005"
var EnableAdditionalCommands = os.Getenv("ENABLE_ADDITIONAL_COMMANDS") == "true"
var ItToolsUrl = os.Getenv("IT_TOOLS_URL")
var ScriptsPath = "user_scripts"
var HostUrl = os.Getenv("HOST_URL")
