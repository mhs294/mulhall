package main

import (
	"log"

	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/server"
)

var app *server.AppServer

func init() {
	log.Println("init started")

	err := env.LoadVars()
	if err != nil {
		panic(err)
	}
	app, err = server.NewAppServer()
	if err != nil {
		panic(err)
	}

	log.Println("init complete")
}

func main() {
	app.Start()
}
