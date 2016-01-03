package api

import (
	"encoding/json"
	"go-lbapp/db"
	"go-lbapp/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// UserHandler : Return function for handler
func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]

	if userid == "" {
		log.Fatal("No user Id supplied")
		ThrowInternalErrAndExit(w)
	}

	// Extract a user from DB
	var u models.User

	// Used for per user connection to DB
	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	err := dbconn.
		QueryRow("SELECT userid, firstname, lastname, gender, age, phonenumber FROM users WHERE userid = $1", userid).
		Scan(&u.UserID, &u.FirstName, &u.LastName, &u.Gender, &u.Age, &u.PhoneNumber)

	if err != nil {
		// If execution err occurs then throw error
		log.Fatal(err)
		ThrowForbiddenedAndExit(w)
	}

	jsonResponse, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append the data to response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
