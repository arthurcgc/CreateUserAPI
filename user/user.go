package user

type User struct {
	name  string
	email string
}

func (u *User) GetUserName() string {
	return u.name
}

func (u *User) GetUserEmail() string {
	return u.email
}
