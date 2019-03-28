package main

import (
	"log"

	"github.com/nsogame/api"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var err error

	config, err := api.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := api.NewInstance(&config)
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
}
