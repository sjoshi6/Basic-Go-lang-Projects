package api

import (
	"database/sql"
	"encoding/json"
	"expvar"
	"fmt"
	"go-lbapp/config"
	"go-lbapp/db"
	"go-lbapp/generics"
	"go-lbapp/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (

	//DBName : Used for conenctions to database
	DBName              = settings.DBName
	eventsTableName     = settings.EventsTableName
	cost            int = settings.Cost
)

// Map for number of route hits
var routeHits = expvar.NewMap("routeHits").Init()

/* Contains all the Route Handlers for API function calls */

// CreateAccount : Handle Signup requests for new user
func CreateAccount(w http.ResponseWriter, r *http.Request) {

	routeHits.Add("/v1/signup", 1)

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
	stmt, _ := dbconn.Prepare("INSERT INTO Users(UserId,Password,FirstName,LastName,Gender,Age,PhoneNumber) VALUES($1,$2,$3,$4,$5,$6,$7);")

	hash, err := bcrypt.GenerateFromPassword([]byte(signupdata.Password), cost)
	if err != nil {

		fmt.Println("bcrypt hash creation broke")
		ThrowInternalErrAndExit(w)

	} else {

		age, _ := strconv.ParseInt(signupdata.Age, 10, 64)

		_, err := stmt.Exec(signupdata.UserID,
			string(hash),
			signupdata.FirstName,
			signupdata.LastName,
			signupdata.Gender,
			age,
			signupdata.PhoneNumber)

		if err != nil {
			log.Fatal(err)
		}

		RespondSuccessAndExit(w, "User Registered Successfully")
	}
}

// ConfirmCredentials : Handle Login requests for existing users
func ConfirmCredentials(w http.ResponseWriter, r *http.Request) {

	routeHits.Add("/v1/login", 1)

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

	rows, err := dbconn.Query("SELECT Password FROM Users where UserID='" + string(logindata.UserID) + "'")
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

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Status", "Client Error")
		RespondOrThrowErr(responsecontent, w)
		return
	}

	// If no error in comparehash means login Credentials match
	RespondSuccessAndExit(w, "User Login Successful")

}

// SearchEventsByRange : Used to search events created in a chosen radius
func SearchEventsByRange(w http.ResponseWriter, r *http.Request) {

	// for unpacking events
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
	)

	var returnEvents generics.Events
	var searchevents generics.SearchEventsData

	decoder := json.NewDecoder(r.Body)
	// Expand the json attached in post request

	err := decoder.Decode(&searchevents)
	if err != nil {
		panic(err)
	}

	// Create a geo point using lat & longitude
	lat, _ := strconv.ParseFloat(searchevents.Lat, 64)
	long, _ := strconv.ParseFloat(searchevents.Long, 64)
	radius, _ := strconv.ParseFloat(searchevents.Radius, 64)

	events, err := getEnventsByRange(lat, long, radius)

	if err != nil {
		fmt.Println(err)
	}

	for events.Next() {
		err := events.Scan(&id, &eventname, &lat, &long,
			&creationtime, &creatorid, &starttime, &endtime, &maxMem, &minMem,
			&friendOnly, &gender, &minAge, &maxAge)

		if err != nil {
			log.Fatal(err)
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

		returnEvents = append(returnEvents, event)

	}

	// Create a JSON to reply to the client
	reply := generics.SearchResults{returnEvents}
	jsonReply, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append the data to response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonReply)
}

func getEnventsByRange(lat, long, radius float64) (*sql.Rows, error) {

	selectStr := fmt.Sprintf("SELECT * FROM %v a", eventsTableName)
	lat1 := fmt.Sprintf("sin(radians(%f)) * sin(radians(a.lat))", lat)
	lng1 := fmt.Sprintf("cos(radians(%f)) * cos(radians(a.lat)) * cos(radians(a.lng) - radians(%f))", lat, long)
	whereStr := fmt.Sprintf("WHERE acos(%s + %s) * %f <= %f", lat1, lng1, 6356.7523, radius)
	query := fmt.Sprintf("%s %s", selectStr, whereStr)

	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	res, err := dbconn.Query(query)
	if err != nil {
		panic(err)
	}

	return res, err
}
