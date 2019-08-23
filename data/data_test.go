package data

import (
	"database/sql"
	"testing"

	"github.com/arthurcgc/CreateUserAPI/myuser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db = data{username: "root", password: "root", database: nil}

func TestOpenDb(t *testing.T) {
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()

	require.NotEqual(t, db.database, nil, "Error opening database\n")
}

func TestInsertUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotEqual(t, db.database, nil, "Error opening database\n")

	defer db.CloseDb()
	t.Run("", func(t *testing.T) {
		err := db.InsertUser("Arthur", "arthur@gmail.com")
		assert.NoError(t, err)

		var rows *sql.Rows
		rows, err = db.database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
		require.Equal(t, err, nil, err)

		assert.True(t, rows.Next())
		var name string
		var email string
		err = rows.Scan(&name, &email)
		assert.NoError(t, err)
		assert.Equal(t, name, "Arthur", "name differs\n")
		assert.Equal(t, email, "arthur@gmail.com", "email differs\n")
		assert.False(t, rows.Next())
		// assert.Equal(t, rows.Next(), false, "two or more rows exist\n")
	})

	cleanUpDatabase(t)

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		assert.NoError(t, err)

		err = db.InsertUser("Arthur", "arthur@gmail.com")
		assert.Error(t, err)
		// assert.NotEqual(t, err, nil, "User inserted and should not be inserted\n")
	})

	// teardown
	cleanUpDatabase(t)
}

func cleanUpDatabase(t *testing.T) {
	err := db.OpenDb()
	require.NoError(t, err)

	_, err = db.database.Exec("DELETE FROM User")
	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotEqual(t, db.database, nil, "Error opening database\n")

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		assert.NoError(t, err)

		err = db.DeleteUser("Arthur", "arthur@gmail.com")
		assert.NoError(t, err)
		var rows *sql.Rows
		rows, err = db.database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
		assert.NoError(t, err)

		assert.False(t, rows.Next())
	})

	// teardown
	cleanUpDatabase(t)
}

func TestGetUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotEqual(t, db.database, nil, "Error opening database\n")

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		assert.NoError(t, err)

		var got *myuser.User
		got, err = db.GetUser("arthur@gmail.com")
		assert.NoError(t, err)
		assert.NotEqual(t, got, nil, "User returned nil")

		assert.Equal(t, got.Email, "arthur@gmail.com", "Error querying email")
	})

	t.Run("", func(t *testing.T) {
		// var got *user.User
		_, err = db.GetUser("notfound@gmail.com")
		assert.Error(t, err)
	})

	// teardown
	cleanUpDatabase(t)
}
