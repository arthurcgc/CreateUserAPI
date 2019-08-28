package data

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func createMockDb() (*Data, sqlmock.Sqlmock, error) {
	data := new(Data)
	DB, mock, err := sqlmock.New()

	data.Database = DB
	data.Username = "root"
	data.Password = "root"
	return data, mock, err
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetOne(t *testing.T) {
	data, mock, err := createMockDb()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// columns := []string{"o_Name", "o_Email"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`"SELECT * from User WHERE Email=?"`)).WithArgs("arthur@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"Name", "Email"}).AddRow("Arthur", "arthur@gmail.com"))
	mock.ExpectCommit()

	data.Database.Begin()
	_, err = data.GetUser("arthur@gmail.com")
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert(t *testing.T) {
	data, mock, err := createMockDb()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// columns := []string{"o_Name", "o_Email"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`"INSERT INTO "User" VALUES (?, ?);"`)).WithArgs("Arthur", "arthur@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"Name", "Email"}).AddRow(1, 1))
	mock.ExpectCommit()

	data.Database.Begin()
	err = data.InsertUser("Arthur", "arthur@gmail.com")
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
