package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/internal/server"
)

func main() {
	if err := server.StartServer(config.DbPath, config.Port); err != nil {
		panic(err)
	}
}
