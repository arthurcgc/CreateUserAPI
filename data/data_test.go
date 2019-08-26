package data

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db = data{username: "root", password: "root", database: nil}

func TestOpenDb(t *testing.T) {
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotNil(t, db.database)
}

func TestInsertUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotNil(t, db.database)

	defer db.CloseDb()
	t.Run("", func(t *testing.T) {
		err := db.InsertUser("Arthur", "arthur@gmail.com")
		assert.NoError(t, err)

		var rows *sql.Rows
		rows, err = db.database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
		require.NoError(t, err)

		assert.True(t, rows.Next())
		var name string
		var email string
		err = rows.Scan(&name, &email)
		assert.NoError(t, err)
		assert.Equal(t, name, "Arthur", "name differs\n")
		assert.Equal(t, email, "arthur@gmail.com", "email differs\n")
		assert.False(t, rows.Next())
	})

	cleanUpDatabase(t)

	t.Run("", func(t *testing.T) {
		_, err := db.database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		assert.NoError(t, err)

		err = db.InsertUser("Arthur", "arthur@gmail.com")
		assert.Error(t, err)
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
	require.NotNil(t, db.database)

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
	require.NotNil(t, db.database)

	t.Run("", func(t *testing.T) {
		expectedUser := &myuser.User{Name: "Arthur", Email: "arthur@gmail.com"}
		_, err := db.database.Exec("INSERT INTO User VALUES (?, ?)", expectedUser.Name, expectedUser.Email)
		assert.NoError(t, err)

		var got *myuser.User
		got, err = db.GetUser("arthur@gmail.com")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got, expectedUser, "Error querying email")
	})

	t.Run("", func(t *testing.T) {
		// var got *user.User
		_, err = db.GetUser("notfound@gmail.com")
		assert.Error(t, err)
	})

	// teardown
	cleanUpDatabase(t)
}
