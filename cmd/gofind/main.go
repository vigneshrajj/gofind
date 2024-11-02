package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/internal/server"
)

func main() {
	server.StartServer()
}
