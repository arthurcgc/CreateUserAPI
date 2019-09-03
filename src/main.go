package main

import (
	"github.com/arthurcgc/CreateUserAPI/restapi"
)

// func getDbCredentials(in *os.File) (string, string, error) {
// 	if in == nil {
// 		in = os.Stdin
// 	}
// 	// fmt.Printf("username: ")
// 	// _, err := fmt.Fscanf(in, "%s", &username)
// 	username := os.Getenv("SQL_USER")
// 	pswd := os.Getenv("SQL_PSSWD")

// 	// if err != nil {
// 	// 	return "", "", err
// 	// }
// 	// fmt.Println("Your password: ")
// 	// bytePassword, _ := terminal.ReadPassword(int(in.Fd()))
// 	// fmt.Println() // it's necessary to add a new line after user's input

// 	return username, pswd, nil
// }

func main() {
	// username, pswd, err := getDbCredentials(os.Stdin)
	// if err != nil {
	// 	panic(err)
	// }
	app, err := restapi.Initialize()
	if err != nil {
		panic(err)
	}

	app.Run()
}
