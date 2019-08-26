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
	Database data.Data
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
	// json.NewEncoder(w).Encode(usr)
	respondWithJSON(w, http.StatusOK, usr)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request, db *data.Data) {

}

func InsertUser(w http.ResponseWriter, r *http.Request, db *data.Data) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request, db *data.Data) {

}

func (app *RestApi) appendRouterFunctions() {
	app.Router.HandleFunc("/users/{email}", app.GetUser).Methods("GET")
}

func Initialize() (*RestApi, error) {
	app := new(RestApi)
	app.Database.Username = "root"
	app.Database.Password = "root"
	err := app.Database.OpenDb()
	if err != nil {
		return nil, err
	}
	app.Router = mux.NewRouter()
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
