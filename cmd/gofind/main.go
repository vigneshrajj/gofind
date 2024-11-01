package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
)


func main() {
	config.StartServer()
}
