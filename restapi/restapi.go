package restapi

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arthurcgc/CreateUserAPI/myuser"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *data.data) myuser.User {
	db.OpenDb()
	defer db.CloseDb()
	usr := db.GetUser(email)
	return usr
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
