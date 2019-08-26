package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/arthurcgc/CreateUserAPI/restapi"
)

var app *restapi.RestApi

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

// func TestMain(m *testing.M) {
// 	var err error
// 	app, err = restapi.Initialize()
// 	require.NotNil(m, err)

// 	code := m.Run()

// 	clearTable()

// 	os.Exit(code)
// }
