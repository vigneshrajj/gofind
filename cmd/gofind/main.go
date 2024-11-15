package main

import (
	"github.com/vigneshrajj/gofind/config"
	"github.com/vigneshrajj/gofind/internal/server"
)

func main() {
	if err := server.StartServer(config.DbPath, config.Port); err != nil {
		panic(err)
	}
}
