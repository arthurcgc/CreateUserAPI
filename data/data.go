package data

import (
	"database/sql"
	"fmt"
	"os"

	//comment justifying it
	"github.com/arthurcgc/CreateUserAPI/myuser"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh/terminal"
)

type data struct {
	username string
	password string
	database *sql.DB
}

func (db *data) SetDbCredentials(in *os.File) error {
	if in == nil {
		in = os.Stdin
	}
	var username string
	fmt.Printf("username: ")
	_, err := fmt.Fscanf(in, "%s", &username)
	if err != nil {
		return err
	}
	fmt.Println("Your password: ")
	bytePassword, _ := terminal.ReadPassword(int(in.Fd()))
	fmt.Println() // it's necessary to add a new line after user's input
	db.username = username
	db.password = string(bytePassword)
	return nil
}

func (db *data) OpenDb() error {
	var err error
	db.database, err = sql.Open("mysql", db.getDbConnectionString())
	if err != nil {
		return err
	}
	return nil
}

func (db *data) getDbConnectionString() string {
	dbString := db.username + ":" + db.password + "@/second_go_proj"
	return dbString
}

func (db *data) CloseDb() {
	db.database.Close()
}

func (db *data) InsertUser(name string, email string) error {
	stmtIns, err := db.database.Prepare("INSERT INTO User VALUES (?, ?);")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(name, email)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func (db *data) DeleteUser(name string, email string) error {
	stmtRm, err := db.database.Prepare("DELETE FROM User WHERE Name = ? AND Email = ?")
	if err != nil {
		return err
	}
	_, err = stmtRm.Exec(name, email)
	if err != nil {
		return err
	}
	defer stmtRm.Close()
	return nil
}

func (db *data) GetUser(email string) (*myuser.User, error) {
	rows, err := db.database.Query("SELECT * from User WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("User Not Found")
	}
	var userName, userEmail string
	err = rows.Scan(&userName, &userEmail)
	if err != nil {
		return nil, err
	}

	res := &myuser.User{Name: userName, Email: userEmail}
	return res, nil
}
