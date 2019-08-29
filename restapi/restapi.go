package restapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arthurcgc/CreateUserAPI/data"
	"github.com/gorilla/mux"
)

type RestApi struct {
	Router   *mux.Router
	Database data.DataInterface
}

func (app *RestApi) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error opening Database")
		return
	}
	defer app.Database.CloseDb()

	usr, err := app.Database.GetUser(params["email"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error retrieving User")
		return
	}
	if usr != nil {
		respondWithJSON(w, http.StatusOK, usr)
		return
	}
	respondWithJSON(w, http.StatusNotFound, usr)
}

func (app *RestApi) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error opening Database")
		return
	}
	defer app.Database.CloseDb()

	usrs, err := app.Database.GetAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error retrieving Users")
		return
	}
	// json.NewEncoder(w).Encode(usr)
	if len(usrs) == 0 {
		respondWithJSON(w, http.StatusNotFound, usrs)
	}
	respondWithJSON(w, http.StatusOK, usrs)
}

func (app *RestApi) InsertUser(w http.ResponseWriter, r *http.Request) {
	helper := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&helper); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if len(helper.Name) == 0 || len(helper.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error opening Database")
		return
	}
	defer app.Database.CloseDb()
	if err := app.Database.InsertUser(helper.Name, helper.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error inserting user to Database")
		return
	}
	str := "Name: " + helper.Name + " Email: " + helper.Email
	respondWithJSON(w, http.StatusCreated, str)
}

func (app *RestApi) UpdateUser(w http.ResponseWriter, r *http.Request) {
	helper := struct {
		Email    string `json:"email"`
		NewName  string `json:"newName"`
		NewEmail string `json:"newEmail"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&helper); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error opening Database")
		return
	}
	defer app.Database.CloseDb()
	usr, err := app.Database.UpdateUser(helper.Email, helper.NewEmail, helper.NewName)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error updating user to Database")
		return
	}
	str := "User: " + usr.Name + ", with email: " + usr.Email + " updated"
	respondWithJSON(w, http.StatusOK, str)
}

func (app *RestApi) DeleteUser(w http.ResponseWriter, r *http.Request) {
	helper := struct {
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&helper); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error opening Database")
		return
	}
	defer app.Database.CloseDb()
	usr, err := app.Database.DeleteUser(helper.Email)
	if err != nil || usr == nil {
		respondWithError(w, http.StatusBadRequest, "Error deleting user from Database")
		return
	}
	str := "User: " + usr.Name + ", with email: " + usr.Email + " deleted"
	respondWithJSON(w, http.StatusOK, str)
}

func (app *RestApi) appendRouterFunctions() {
	app.Router.HandleFunc("/users/{email}", app.GetUser).Methods("GET")
	app.Router.HandleFunc("/users", app.GetAllUsers).Methods("GET")
	app.Router.HandleFunc("/users/", app.InsertUser).Methods("POST")
	app.Router.HandleFunc("/users/", app.UpdateUser).Methods("UPDATE")
	app.Router.HandleFunc("/users/", app.DeleteUser).Methods("DELETE")
}

func Initialize(username string, password string) (*RestApi, error) {
	db := &data.Data{Username: username, Password: password}
	// app.Database.Username = username
	// app.Database.Password = password

	err := db.OpenDb()
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	app := &RestApi{Router: r, Database: db}

	app.appendRouterFunctions()
	return app, nil
}

func (app *RestApi) Run() {
	log.Fatal(http.ListenAndServe(":8000", app.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
