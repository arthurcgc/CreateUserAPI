package restapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthurcgc/CreateUserAPI/data"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeGetRequest(t *testing.T, handler http.HandlerFunc, path string) *http.Response {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()
	client := ts.Client()
	res, err := client.Get(ts.URL + path)
	require.NoError(t, err)
	return res
}

func initApp(mock *mockData) *RestApi {
	r := mux.NewRouter()
	app := &RestApi{Router: r, Database: mock}
	return app
}

func TestGetAll(t *testing.T) {

	t.Run("", func(t *testing.T) {
		// setup
		var users []data.User
		users = append(users, data.User{Name: "test1", Email: "test1@gmail.com"})
		users = append(users, data.User{Name: "test2", Email: "test2@gmail.com"})
		mock := &mockData{
			getAllfunc: func() ([]data.User, error) {
				return users, nil
			}}
		app := initApp(mock)
		res := executeGetRequest(t, app.GetAllUsers, "/users")
		require.Equal(t, http.StatusOK, res.StatusCode)
		info, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t,
			"[{\"Name\":\"test1\",\"Email\":\"test1@gmail.com\"},{\"Name\":\"test2\",\"Email\":\"test2@gmail.com\"}]",
			string(info))
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &(mockData{})

		app := initApp(mock)

		res := executeGetRequest(t, app.GetAllUsers, "/users")
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		mock.getAllfunc = func() ([]data.User, error) {
			return nil, errors.New("Forced Error")
		}
		app := initApp(mock)
		res := executeGetRequest(t, app.GetAllUsers, "/users/")
		require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// setup
		usr := &data.User{Name: "test", Email: "test@gmail.com"}
		mock := &mockData{getUserfunc: func(email string) (*data.User, error) {
			return usr, nil
		}}

		app := initApp(mock)
		res := executeGetRequest(t, app.GetUser, "/users/test@gmail.com")
		require.Equal(t, http.StatusOK, res.StatusCode)

		assertReliableData(t, res, usr)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}

		app := initApp(mock)
		res := executeGetRequest(t, app.GetUser, "/users/test@gmail.com")
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		mock.getUserfunc = func(email string) (*data.User, error) {
			return nil, errors.New("Forced Error")
		}
		app := initApp(mock)
		res := executeGetRequest(t, app.GetUser, "/users/test@gmail.com")
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func executePostRequest(t *testing.T, app *RestApi, path string, payload []byte) *http.Response {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InsertUser(w, r)
	}))
	defer ts.Close()
	client := ts.Client()
	res, err := client.Post(ts.URL+path, "application/json", bytes.NewBuffer(payload))
	require.NoError(t, err)
	return res
}

func assertReliableData(t *testing.T, res *http.Response, usr *data.User) {
	response := data.User{}
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err)
	require.Equal(t, usr.Name, response.Name)
	require.Equal(t, usr.Email, response.Email)
}

func setUserPayload(t *testing.T, usr *data.User) []byte {
	payload, err := json.Marshal(usr)
	require.NoError(t, err)
	return payload
}

func TestInsertUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// setup
		usr := &data.User{Name: "test", Email: "test@gmail.com"}
		mock := &mockData{}
		mock.insertUserfunc = func(name, email string) (*data.User, error) {
			return usr, nil
		}
		app := initApp(mock)

		payload := setUserPayload(t, usr)
		res := executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		assertReliableData(t, res, usr)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executePostRequest(t, app, "/users/", nil)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		payload, err := json.Marshal(data.User{Name: "", Email: ""})
		require.NoError(t, err)
		res := executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusNotAcceptable, res.StatusCode)

		payload, err = json.Marshal(data.User{Name: "test", Email: ""})
		require.NoError(t, err)
		res = executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusNotAcceptable, res.StatusCode)

		payload, err = json.Marshal(data.User{Name: "", Email: "test@gmail.com"})
		require.NoError(t, err)
		res = executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusNotAcceptable, res.StatusCode)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		mock.insertUserfunc = func(name, email string) (*data.User, error) {
			return nil, errors.New("Forced error")
		}
		app := initApp(mock)
		payload := setUserPayload(t, &data.User{Name: "test", Email: "test@gmail.com"})

		res := executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusNotAcceptable, res.StatusCode)
	})
}

func executeDeleteRequest(t *testing.T, app *RestApi, path string, payload []byte) *http.Response {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.DeleteUser(w, r)
	}))
	defer ts.Close()
	client := ts.Client()
	req, err := http.NewRequest("DELETE", ts.URL+path, bytes.NewBuffer(payload))
	require.NoError(t, err)
	res, err := client.Do(req)
	assert.NoError(t, err)
	return res
}

func TestDeleteUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// setup
		usr := &data.User{Email: "test@gmail.com"}
		mock := &mockData{deleteUserfunc: func(email string) (*data.User, error) {
			return &data.User{Email: email}, nil
		}}
		app := initApp(mock)
		payload := setUserPayload(t, usr)
		res := executeDeleteRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusAccepted, res.StatusCode)
		assertReliableData(t, res, usr)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		app := initApp(mock)
		res := executeDeleteRequest(t, app, "/users/", nil)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)

		payload := setUserPayload(t, &data.User{Name: "test", Email: "test@gmail.com"})
		res = executeDeleteRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		mock.deleteUserfunc = func(email string) (*data.User, error) {
			return nil, errors.New("Forced Error")
		}
		app := initApp(mock)

		payload := setUserPayload(t, &data.User{Name: "test", Email: "test@gmail.com"})
		res := executeDeleteRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func executeUpdateRequest(t *testing.T, app *RestApi, path string, payload []byte) *http.Response {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.UpdateUser(w, r)
	}))
	defer ts.Close()
	client := ts.Client()
	req, err := http.NewRequest("UPDATE", ts.URL+path, bytes.NewBuffer(payload))
	require.NoError(t, err)
	res, err := client.Do(req)
	assert.NoError(t, err)
	return res
}

func TestUpdateUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// setup
		helper := updateArgs{
			Email:    "bla@xxx.com",
			NewName:  "blu",
			NewEmail: "blu@xxx.com",
		}
		mock := &mockData{}
		mock.updateUserfunc = func(email, newEmail, newName string) (*data.User, error) {
			return &data.User{Name: helper.NewName, Email: helper.NewEmail}, nil
		}
		app := initApp(mock)
		payload, err := json.Marshal(helper)
		require.NoError(t, err)
		res := executeUpdateRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)

		usr := &data.User{Name: helper.NewName, Email: helper.NewEmail}
		assertReliableData(t, res, usr)
	})
}
