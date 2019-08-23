package data

import (
	"database/sql"
	"testing"

	"../user"
)

var db = data{username: "root", password: "root", database: nil}

func TestOpenDb(t *testing.T) {
	err := db.OpenDb()
	if db.database == nil || err != nil {
		t.Errorf("Error opening database\n")
	}
}

func TestInsertUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	if db.database == nil || err != nil {
		t.Errorf("Error opening database\n")
	}
	defer db.CloseDb()
	t.Run("", func(t *testing.T) {
		err := db.InsertUser("Arthur", "arthur@gmail.com")
		if err != nil {
			t.Errorf("Error inserting user to User table\n")
		}
		var rows *sql.Rows
		rows, err = db.database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
		if err != nil {
			t.Fatalf("%v", err)
		}
		if !rows.Next() {
			t.Fatalf("No row found\n")
		}
		var name string
		var email string
		err = rows.Scan(&name, &email)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if name != "Arthur" {
			t.Fatalf("name differs\n")
		} else if email != "arthur@gmail.com" {
			t.Fatalf("email differs\n")
		}
		if rows.Next() {
			t.Fatalf("two or more rows exist\n")
		}
	})

	cleanUpDatabase(t)

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		if err != nil {
			t.Errorf("Error inserting user\n")
		}
		err = db.InsertUser("Arthur", "arthur@gmail.com")
		if err == nil {
			t.Errorf("User inserted and should not be inserted\n")
		}
	})

	// teardown
	cleanUpDatabase(t)
}

func cleanUpDatabase(t *testing.T) {
	if err := db.OpenDb(); err != nil {
		t.Fatalf("could not open database connection")
	}

	_, err := db.database.Exec("DELETE FROM User")
	if err != nil {
		t.Fatalf("error during database cleanup: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	if db.database == nil || err != nil {
		t.Errorf("Error opening database\n")
	}
	defer db.CloseDb()

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		if err != nil {
			t.Errorf("Error inserting user\n")
		}
		err = db.DeleteUser("Arthur", "arthur@gmail.com")
		if err != nil {
			t.Errorf("Error while deleting user\n%v", err)
		}
		var rows *sql.Rows
		rows, err = db.database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
		if err != nil {
			t.Fatalf("%v", err)
		}
		if rows.Next() {
			t.Fatalf("Row still exists\n")
		}

	})

	// teardown
	cleanUpDatabase(t)
}

func TestGetUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	if db.database == nil || err != nil {
		t.Errorf("Error opening database\n")
	}
	defer db.CloseDb()

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		if err != nil {
			t.Errorf("Error inserting user\n")
		}
		var got *user.User
		got, err = db.GetUser("arthur@gmail.com")
		if err != nil {
			t.Errorf("Error in GetUser: %v", err)
		}
		if got == nil {
			t.Errorf("User returned nil")
		} else if got.Email != "arthur@gmail.com" {
			t.Errorf("Error querying email")
		}
	})

	t.Run("", func(t *testing.T) {
		// var got *user.User
		_, err = db.GetUser("notfound@gmail.com")
		if err == nil {
			t.Errorf("Error should not be nil: %v", err)
		}
	})

	// teardown
	cleanUpDatabase(t)
}
