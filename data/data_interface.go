package data

type DataInterface interface {
	OpenDb() error
	CloseDb() // should return error
	InsertUser(name string, email string) error
	UpdateUser(email string, newEmail string, newName string) (*User, error)
	DeleteUser(email string) (*User, error)
	GetUser(email string) (*User, error)
	GetAll() ([]User, error)
}
