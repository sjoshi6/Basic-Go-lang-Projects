package api

import (
	"encoding/json"
	"fmt"
	"go-lbapp/db"
	"go-lbapp/generics"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const (

	//DBName : Used for conenctions to database
	DBName     = "db_lbapp"
	cost   int = 10
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

	stmt, _ := dbconn.Prepare("INSERT INTO userlogin(UserID,Password,name) VALUES($1,$2,$3);")

	/*  Add a DB insert command here to register new users
	    Ensure new user credentials are hashed with bcrypt
	*/
	hash, err := bcrypt.GenerateFromPassword([]byte(signupdata.Password), cost)
	if err != nil {

		fmt.Println("bcrypt hash creation broke")
		responsecontent := BasicResponse{
			"Internal Server Error",
			500,
		}

		w.Header().Set("StatusCode", "500")
		w.Header().Set("Status", "Internal Server Error")
		respondOrThrowErr(responsecontent, w)

	} else {

		_, err := stmt.Exec(string(signupdata.UserID), string(hash), string(signupdata.Name))
		if err != nil {
			log.Fatal(err)
		}

		responsecontent := BasicResponse{
			"User Registered Successfully",
			200,
		}
		respondOrThrowErr(responsecontent, w)

	}
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
