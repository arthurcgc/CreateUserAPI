package restapi

import (
	"github.com/arthurcgc/CreateUserAPI/data"
)

type mockData struct {
	openDbFunc     func() error
	closeDbfunc    func() error
	insertUserfunc func(name, email string) error
	updateUserfunc func(email, newEmail, newName string) (*data.User, error)
	deleteUserfunc func(email string) (*data.User, error)
	getUserfunc    func(email string) (*data.User, error)
	getAllfunc     func() ([]data.User, error)
}

func (m *mockData) OpenDb() error {
	if m.openDbFunc == nil {
		return nil
	}
	return m.openDbFunc()
}

func (m *mockData) CloseDb() error {
	if m.closeDbfunc == nil {
		return nil
	}
	return m.closeDbfunc()
}

func (m *mockData) InsertUser(name, email string) error {
	if m.insertUserfunc == nil {
		return nil
	}
	return m.insertUserfunc(name, email)
}

func (m *mockData) UpdateUser(email, newEmail, newName string) (*data.User, error) {
	if m.updateUserfunc == nil {
		return nil, nil
	}
	return m.updateUserfunc(email, newEmail, newName)
}

func (m *mockData) DeleteUser(email string) (*data.User, error) {
	if m.deleteUserfunc == nil {
		return nil, nil
	}
	return m.deleteUserfunc(email)
}

func (m *mockData) GetUser(email string) (*data.User, error) {
	if m.getUserfunc == nil {
		return nil, nil
	}
	return m.getUserfunc(email)
}

func (m *mockData) GetAll() ([]data.User, error) {
	if m.getAllfunc == nil {
		return nil, nil
	}
	return m.getAllfunc()
}
