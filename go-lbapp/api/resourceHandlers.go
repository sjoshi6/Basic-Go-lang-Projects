package api

import (
	"encoding/json"
	"fmt"
	"go-lbapp/db"
	"go-lbapp/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// UserHandler : Get User
func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]

	if userid == "" {
		ThrowForbiddenedAndExit(w)
		return
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

		ThrowForbiddenedAndExit(w)
		return
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

// EventHandler : Get Event
func EventHandler(w http.ResponseWriter, r *http.Request) {

	// Vars to handle Event Struct
	var (
		id           string
		eventname    string
		creationtime time.Time
		creatorid    string
		starttime    time.Time
		endtime      time.Time
		maxMem       int64
		minMem       int64
		friendOnly   bool
		gender       string
		minAge       int64
		maxAge       int64
		lat          float64
		long         float64
	)

	// Extract API resource ID
	vars := mux.Vars(r)
	eventid := vars["eventid"]

	if eventid == "" {
		ThrowForbiddenedAndExit(w)
		return
	}

	// Used for per user connection to DB
	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	err := dbconn.
		QueryRow("SELECT id, event_name, lat, lng, creation_time, creator_id, start_time, end_time, max_mem, min_mem, friend_only, gender, min_age, max_age FROM Events WHERE id = $1", eventid).
		Scan(&id, &eventname, &lat, &long,
		&creationtime, &creatorid, &starttime, &endtime, &maxMem, &minMem,
		&friendOnly, &gender, &minAge, &maxAge)

	if err != nil {
		// If execution err occurs then throw error
		fmt.Println(err)
		ThrowForbiddenedAndExit(w)
		return
	}

	event := models.Event{

		id,
		eventname,
		strconv.FormatFloat(lat, 'f', 6, 64),
		strconv.FormatFloat(long, 'f', 6, 64),
		creationtime.Format("2014-06-08T02:02:22Z"),
		creatorid,
		starttime.Format("2014-06-08T02:02:22Z"),
		endtime.Format("2014-06-08T02:02:22Z"),
		strconv.FormatInt(maxMem, 10),
		strconv.FormatInt(minMem, 10),
		strconv.FormatBool(friendOnly),
		gender,
		strconv.FormatInt(minAge, 10),
		strconv.FormatInt(maxAge, 10),
	}

	jsonResponse, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append the data to response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
