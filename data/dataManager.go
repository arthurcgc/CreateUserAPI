package data

type User struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

type DataManager interface {
	OpenDb() error
	CloseDb() error
	InsertUser(name string, email string) (*User, error)
	UpdateUser(email string, newEmail string, newName string) (*User, error)
	DeleteUser(email string) (*User, error)
	GetUser(email string) (*User, error)
	GetAll() ([]User, error)
}
