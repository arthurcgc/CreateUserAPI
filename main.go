package main

import "github.com/arthurcgc/CreateUserAPI/restapi"

func main() {
	app, err := restapi.Initialize()
	if err != nil {
		panic(err)
	}
	app.Run()
}
