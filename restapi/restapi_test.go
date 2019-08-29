package restapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthurcgc/CreateUserAPI/data"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEmptyMock() *mockData {
	mock := &mockData{openDbFunc: nil, closeDbfunc: nil, insertUserfunc: nil,
		updateUserfunc: nil, deleteUserfunc: nil, getUserfunc: nil,
		getAllfunc: nil}

	return mock
}

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

func TestGetAll(t *testing.T) {

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{
			getAllfunc: func() ([]data.User, error) {
				var users []data.User
				users = append(users, data.User{Name: "test1", Email: "test1@gmail.com"})
				users = append(users, data.User{Name: "test2", Email: "test2@gmail.com"})
				return users, nil
			}}

		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executeGetRequest(t, app.GetAllUsers, "/users")
		require.Equal(t, http.StatusOK, res.StatusCode)
		info, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		// [{"Name":"test1","Email":"test1@gmail.com"},{"Name":"test2","Email":"test2@gmail.com"}]
		assert.Equal(t,
			"[{\"Name\":\"test1\",\"Email\":\"test1@gmail.com\"},{\"Name\":\"test2\",\"Email\":\"test2@gmail.com\"}]",
			string(info))
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &(mockData{})

		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executeGetRequest(t, app.GetAllUsers, "/users")
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{getUserfunc: func(email string) (*data.User, error) {
			return &data.User{Name: "test", Email: "test@gmail.com"}, nil
		}}

		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executeGetRequest(t, app.GetUser, "/users/test@gmail.com")
		require.Equal(t, http.StatusOK, res.StatusCode)
		info, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t,
			"{\"Name\":\"test\",\"Email\":\"test@gmail.com\"}",
			string(info))
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}

		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executeGetRequest(t, app.GetUser, "/users/test@gmail.com")
		require.Equal(t, http.StatusNotFound, res.StatusCode)
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

func TestInsertUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		payload, err := json.Marshal(data.User{Name: "test", Email: "test@gmail.com"})
		require.NoError(t, err)
		res := executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		require.NoError(t, err)
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
		require.Equal(t, http.StatusBadRequest, res.StatusCode)

		payload, err = json.Marshal(data.User{Name: "", Email: ""})
		require.NoError(t, err)
		res = executePostRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
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
		mock := &mockData{deleteUserfunc: func(email string) (*data.User, error) {
			return &data.User{Name: "test", Email: "test@email.com"}, nil
		}}
		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		payload, err := json.Marshal(data.User{Name: "test", Email: "test@gmail.com"})
		require.NoError(t, err)
		res := executeDeleteRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{}
		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		// payload, err := json.Marshal(data.User{Name: "test", Email: "test@gmail.com"})
		// require.NoError(t, err)
		res := executeDeleteRequest(t, app, "/users/", nil)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)

		payload, err := json.Marshal(data.User{Name: "test", Email: "test@gmail.com"})
		require.NoError(t, err)
		res = executeDeleteRequest(t, app, "/users/", payload)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
