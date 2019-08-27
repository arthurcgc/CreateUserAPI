package data

type DataInterface interface {
	OpenDb() error
	CloseDb() // should return error
	InsertUser(name string, email string)
	UpdateUser(email string, newEmail string, newName string) (*User, error)
	DeleteUser(email string) (*User, error)
	GetUser(email string)
	GetAll() ([]User, error)
}
