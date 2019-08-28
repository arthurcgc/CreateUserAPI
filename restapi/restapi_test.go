package restapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthurcgc/CreateUserAPI/data"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func setEmptyMock() *mockData {
	mock := &mockData{openDbFunc: nil, closeDbfunc: nil, insertUserfunc: nil,
		updateUserfunc: nil, deleteUserfunc: nil, getUserfunc: nil,
		getAllfunc: nil}

	return mock
}

func executeRequest(t *testing.T, handler http.HandlerFunc, app *RestApi) *http.Response {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.GetAllUsers(w, r)
	}))
	defer ts.Close()
	client := ts.Client()
	res, err := client.Get(ts.URL + "/users")
	require.NoError(t, err)
	return res
}

func TestGetAll(t *testing.T) {

	t.Run("", func(t *testing.T) {
		// setup
		mock := &mockData{openDbFunc: nil, closeDbfunc: nil, insertUserfunc: nil,
			updateUserfunc: nil, deleteUserfunc: nil, getUserfunc: nil,
			getAllfunc: func() ([]data.User, error) {
				var users []data.User
				users = append(users, data.User{Name: "test1", Email: "test1@gmail.com"})
				users = append(users, data.User{Name: "test2", Email: "test2@gmail.com"})
				return users, nil
			}}

		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executeRequest(t, app.GetAllUsers, app)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("", func(t *testing.T) {
		// setup
		mock := &(mockData{})

		r := mux.NewRouter()
		app := &RestApi{Router: r, Database: mock}
		res := executeRequest(t, app.GetAllUsers, app)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
