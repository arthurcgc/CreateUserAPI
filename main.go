package main

import (
	"fmt"
	"os"

	"github.com/arthurcgc/CreateUserAPI/restapi"
	"golang.org/x/crypto/ssh/terminal"
)

func getDbCredentials(in *os.File) (string, string, error) {
	if in == nil {
		in = os.Stdin
	}
	var username string
	fmt.Printf("username: ")
	_, err := fmt.Fscanf(in, "%s", &username)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Your password: ")
	bytePassword, _ := terminal.ReadPassword(int(in.Fd()))
	fmt.Println() // it's necessary to add a new line after user's input
	return username, string(bytePassword), nil
}

func main() {
	username, pswd, err := getDbCredentials(os.Stdin)
	if err != nil {
		panic(err)
	}
	app, err := restapi.Initialize(username, pswd)
	if err != nil {
		panic(err)
	}

	app.Run()
}
