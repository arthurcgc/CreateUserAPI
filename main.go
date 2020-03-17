package main

import (
	"log"

	"github.com/arthurcgc/CreateUserAPI/api"
)

func main() {
	app, err := api.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
