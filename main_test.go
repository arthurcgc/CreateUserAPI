package main

import (
	"net/http"
	"net/http/httptest"
	"os"

	restapi "github.com/arthurcgc/CreateUserAPI/api"
)

var app *restapi.RestApi

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func main() {
	app, err = api.Initialize()
	if err != nil {
		return fmt.Fatalf(err)
	}

	code := m.Run()

	os.Exit(code)
}
