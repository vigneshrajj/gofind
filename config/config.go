package config

import "os"


var DbPath = "db/gofind.db"
var Port = ":3005"
var EnableAdditionalCommands = os.Getenv("ENABLE_ADDITIONAL_COMMANDS") == "true"
