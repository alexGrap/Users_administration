package main

import (
	"avito/config"
	fiberApp "avito/internal/app"
	"log"
)

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
