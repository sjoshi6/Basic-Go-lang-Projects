package api

import (
	"encoding/json"
	"go-lbapp/db"
	"go-lbapp/generics"
	"net/http"
)

const (

	//DBName : Used for conenctions to database
	DBName = "db_lbapp"
)

/* Contains all the Route Handlers for API function calls */

// CreateAccount : Handle Signup requests for new user
func CreateAccount(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var signupdata generics.SignUpData

	// Expand the json attached in post request
	err := decoder.Decode(&signupdata)
	if err != nil {
		panic(err)
	}

	// Used for per user connection to DB
	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	/*  Add a DB insert command here to register new users
	    Ensure new user credentials are hashed with bcrypt
	*/

	responsecontent := BasicResponse{
		"User Registered Successfully",
		200,
	}

	respondOrThrowErr(responsecontent, w)
}

// ConfirmCredentials : Handle Login requests for existing users
func ConfirmCredentials(w http.ResponseWriter, r *http.Request) {

}

// Respond to general requests or exit with server err.
func respondOrThrowErr(responseObj BasicResponse, w http.ResponseWriter) {

	responseJSON, err := json.Marshal(responseObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
