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

	// Add an err handler here to ensure a failed signup request is handled
	stmt, _ := dbconn.Prepare("INSERT INTO userlogin(UserID,Password,name) VALUES($1,$2,$3);")

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

	decoder := json.NewDecoder(r.Body)
	var logindata generics.LoginData

	// Expand the json attached in post request
	err := decoder.Decode(&logindata)
	if err != nil {
		panic(err)
	}

	// Used for per user connection to DB
	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	rows, err := dbconn.Query("SELECT Password FROM userlogin where UserID='" + string(logindata.UserID) + "'")
	var password string

	for rows.Next() {
		rows.Scan(&password)
	}

	loginerr := bcrypt.CompareHashAndPassword([]byte(password), []byte(logindata.Password))
	if loginerr != nil {

		// If err is thrown credentials are mismatched
		responsecontent := BasicResponse{
			"Login Credentials are incorrect",
			400,
		}

		w.Header().Set("StatusCode", "400")
		w.Header().Set("Status", "Client Error")
		respondOrThrowErr(responsecontent, w)
		return
	}

	// If no error in comparehash means login Credentials match
	responsecontent := BasicResponse{
		"User Login Successful",
		200,
	}
	w.Header().Set("StatusCode", "200")
	respondOrThrowErr(responsecontent, w)

}

// CreateEvent : creates a new event at a base location
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var eventcreationdata generics.EventCreationData

	// Expand the json attached in post request
	err := decoder.Decode(&eventcreationdata)
	if err != nil {
		panic(err)
	}

	// Used for per user connection to DB
	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	// Add code to manage event creation request

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
