package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arthurcgc/CreateUserAPI/data"
	"github.com/gorilla/mux"
)

type API struct {
	Router   *mux.Router
	Database data.DataInterface
}

func (app *API) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer app.Database.CloseDb()

	usr, err := app.Database.GetUser(params["email"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if usr != nil {
		respondWithJSON(w, http.StatusOK, usr)
		return
	}
	respondWithJSON(w, http.StatusNotFound, usr)
}

func (app *API) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer app.Database.CloseDb()

	usrs, err := app.Database.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, usrs)
}

func (app *API) InsertUser(w http.ResponseWriter, r *http.Request) {
	helper := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&helper); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if len(helper.Name) == 0 || len(helper.Email) == 0 {
		respondWithError(w, http.StatusNotAcceptable, "Can't insert blank data")
		return
	}
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer app.Database.CloseDb()
	usr, err := app.Database.InsertUser(helper.Name, helper.Email)
	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, usr)
}

type updateArgs struct {
	Email    string `json:"email"`
	NewName  string `json:"newName"`
	NewEmail string `json:"newEmail"`
}

func (app *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var helper updateArgs
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&helper); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if len(helper.Email) == 0 {
		respondWithError(w, http.StatusNotAcceptable, "No primary string passed")
		return
	}
	if len(helper.NewName) == 0 && len(helper.NewEmail) == 0 {
		respondWithError(w, http.StatusNotAcceptable, "No new values passed")
		return
	}
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer app.Database.CloseDb()
	usr, err := app.Database.UpdateUser(helper.Email, helper.NewEmail, helper.NewName)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, usr)
}

func (app *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	helper := struct {
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&helper); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if len(helper.Email) == 0 {
		respondWithError(w, http.StatusNotAcceptable, "Can't accept blank user")
	}
	if err := app.Database.OpenDb(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer app.Database.CloseDb()
	usr, err := app.Database.DeleteUser(helper.Email)
	if err != nil || usr == nil {
		respondWithError(w, http.StatusNotFound, "Invalid user passed")
		return
	}
	respondWithJSON(w, http.StatusAccepted, usr)
}

func (app *API) appendRouterFunctions() {
	app.Router.HandleFunc("/users/{email}", app.GetUser).Methods("GET")
	app.Router.HandleFunc("/users", app.GetAllUsers).Methods("GET")
	app.Router.HandleFunc("/users/", app.InsertUser).Methods("POST")
	app.Router.HandleFunc("/users/", app.UpdateUser).Methods("UPDATE")
	app.Router.HandleFunc("/users/", app.DeleteUser).Methods("DELETE")
}

func Initialize() (*API, error) {
	db := &data.Data{}
	err := db.OpenDb()
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	app := &API{Router: r, Database: db}

	app.appendRouterFunctions()
	return app, nil
}

func (app *API) Run() {
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
