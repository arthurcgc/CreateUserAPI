package data

import (
	"database/sql"
	"fmt"

	//comment justifying it
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name  string
	Email string
}

type Data struct {
	Username string
	Password string
	Database *sql.DB
}

func (db *Data) OpenDb() error {
	var err error
	db.Database, err = sql.Open("mysql", db.getDbConnectionString())
	if err != nil {
		return err
	}
	return nil
}

func (db *Data) getDbConnectionString() string {
	dbString := db.Username + ":" + db.Password + "@/second_go_proj"
	return dbString
}

func (db *Data) CloseDb() {
	db.Database.Close()
}

func (db *Data) InsertUser(name string, email string) error {
	stmtIns, err := db.Database.Prepare("INSERT INTO User VALUES (?, ?);")
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

func (db *Data) UpdateUser(email string, newEmail string, newName string) (*User, error) {
	_, err := db.GetUser(email)
	if err != nil {
		return nil, err
	}
	stmtChangeBoth, err := db.Database.Prepare("UPDATE User SET Name = ?, Email = ? WHERE Email = ?")
	stmtChangeName, err := db.Database.Prepare("UPDATE User SET Name = ? WHERE Email = ?")
	stmtChangeEmail, err := db.Database.Prepare("UPDATE User SET Email = ? WHERE Email = ?")
	if err != nil {
		return nil, err
	}

	var changeEmail bool
	var changeName bool

	if len(newEmail) > 0 {
		changeEmail = true
	}
	if len(newName) > 0 {
		changeName = true
	}
	if !changeEmail && !changeName {
		return nil, fmt.Errorf("Can't update to nil values")
	}
	if changeName && changeEmail {
		_, err = stmtChangeBoth.Exec(newName, newEmail, email)
		if err != nil {
			return nil, err
		}
		return db.GetUser(newEmail)
	}
	if changeName {
		_, err = stmtChangeName.Exec(newName, email)
		if err != nil {
			return nil, err
		}
		return db.GetUser(email)
	}
	_, err = stmtChangeEmail.Exec(newEmail, email)
	if err != nil {
		return nil, err
	}
	return db.GetUser(newEmail)
}

func (db *Data) DeleteUser(email string) (*User, error) {
	usr, err := db.GetUser(email)
	if err != nil {
		return nil, err
	}
	stmtRm, err := db.Database.Prepare("DELETE FROM User WHERE Email = ?")
	if err != nil {
		return nil, err
	}
	_, err = stmtRm.Exec(email)
	if err != nil {
		return nil, err
	}
	defer stmtRm.Close()
	return usr, nil
}

func (db *Data) GetUser(email string) (*User, error) {
	rows, err := db.Database.Query("SELECT * from User WHERE Email=?", email)
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

	res := &User{Name: userName, Email: userEmail}
	return res, nil
}

func (db *Data) GetAll() ([]User, error) {
	rows, err := db.Database.Query("SELECT * from User")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
