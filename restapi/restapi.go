package restapi

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arthurcgc/CreateUserAPI/myuser"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *data.data) (*myuser.User, error) {
	err := db.OpenDb()
	if err != nil {
		return nil, err
	}
	defer db.CloseDb()

	params := mux.Vars(r)
	err, usr := db.GetUser(params["email"])
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, db *data.data) {

}

func InsertUser(w http.ResponseWriter, r *http.Request, db *data.data) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request, db *data.data) {

}

func appendRouterFunctions(router *mux.Router) {
	router.HandleFunc("/users/{email}", GetUser).Methods("GET")
}

func InitializeRouter() *mux.Router {
	router := mux.NewRouter()
	appendRouterFunctions(router)

	return router
}
