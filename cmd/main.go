package main

import (
	"avito/config"
	fiberApp "avito/internal/app"
	"log"
)

// @title avito_intern
// @version 1.0
// @description Intern REST-service.
// @in header
// @host localhost:3000
// @BasePath /
func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	dbConnect, err := config.ParseConfig(conf)
	if err != nil {
		log.Panic(err.Error())
	}
	app := fiberApp.InitApp(dbConnect)
	app.AppStart()
}
