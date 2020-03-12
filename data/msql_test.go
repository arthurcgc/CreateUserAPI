package data

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db = Data{Username: "root", Password: "root", Database: nil}

func TestOpenDb(t *testing.T) {
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotNil(t, db.Database)
}

func TestInsertUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotNil(t, db.Database)

	defer db.CloseDb()
	t.Run("", func(t *testing.T) {
		err := db.InsertUser("Arthur", "arthur@gmail.com")
		assert.NoError(t, err)

		var rows *sql.Rows
		rows, err = db.Database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
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
		_, err := db.Database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
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

	_, err = db.Database.Exec("DELETE FROM User")
	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotNil(t, db.Database)

	t.Run("", func(t *testing.T) {
		_, err := db.Database.Exec("INSERT INTO User VALUES ('Arthur','arthur@gmail.com')")
		assert.NoError(t, err)

		_, err = db.DeleteUser("arthur@gmail.com")
		assert.NoError(t, err)
		var rows *sql.Rows
		rows, err = db.Database.Query("SELECT * from User WHERE name=? AND email=?", "Arthur", "arthur@gmail.com")
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
	require.NotNil(t, db.Database)

	t.Run("", func(t *testing.T) {
		expectedUser := &User{Name: "Arthur", Email: "arthur@gmail.com"}
		_, err := db.Database.Exec("INSERT INTO User VALUES (?, ?)", expectedUser.Name, expectedUser.Email)
		assert.NoError(t, err)

		var got *User
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

func TestGetAll(t *testing.T) {
	// setup
	err := db.OpenDb()
	require.NoError(t, err)
	defer db.CloseDb()
	require.NotNil(t, db.Database)

	var expectedUsers []User
	t.Run("", func(t *testing.T) {
		expectedUsers = append(expectedUsers, User{Name: "Arthur", Email: "arthur@gmail.com"})
		_, err := db.Database.Exec("INSERT INTO User VALUES (?, ?)", "Arthur", "arthur@gmail.com")
		assert.NoError(t, err)

		expectedUsers = append(expectedUsers, User{Name: "Bernardo", Email: "bernardo@gmail.com"})
		_, err = db.Database.Exec("INSERT INTO User VALUES (?, ?)", "Bernardo", "bernardo@gmail.com")
		assert.NoError(t, err)

		expectedUsers = append(expectedUsers, User{Name: "Claudio", Email: "claudio@gmail.com"})
		_, err = db.Database.Exec("INSERT INTO User VALUES (?, ?)", "Claudio", "claudio@gmail.com")
		assert.NoError(t, err)

	})

	t.Run("", func(t *testing.T) {
		// var got *user.User
		users, err := db.GetAll()
		assert.NoError(t, err)
		for i := 0; i < len(users); i++ {
			assert.Equal(t, users[i], expectedUsers[i], "Error while comparing Users")
		}
	})

	// teardown
	cleanUpDatabase(t)
}
