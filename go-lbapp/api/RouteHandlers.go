package api

import (
	"database/sql"
	"encoding/json"
	"expvar"
	"fmt"
	"go-lbapp/db"
	"go-lbapp/generics"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (

	//DBName : Used for conenctions to database
	DBName              = "db_lbapp"
	eventsTableName     = "Events"
	cost            int = 10
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
	stmt, _ := dbconn.Prepare("INSERT INTO userlogin(UserID,Password,name) VALUES($1,$2,$3);")

	hash, err := bcrypt.GenerateFromPassword([]byte(signupdata.Password), cost)
	if err != nil {

		fmt.Println("bcrypt hash creation broke")
		ThrowInternalErrAndExit(w)

	} else {

		_, err := stmt.Exec(string(signupdata.UserID), string(hash), string(signupdata.Name))
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

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Status", "Client Error")
		RespondOrThrowErr(responsecontent, w)
		return
	}

	// If no error in comparehash means login Credentials match
	RespondSuccessAndExit(w, "User Login Successful")

}

// CreateEvent : creates a new event at a base location
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	routeHits.Add("/v1/create_event", 1)

	decoder := json.NewDecoder(r.Body)
	var eventcreationdata generics.EventCreationData

	// Expand the json attached in post request
	err := decoder.Decode(&eventcreationdata)
	if err != nil {
		panic(err)
	}

	// Convert Str input data to respective float / time fmt.
	creationTimeStr := time.Now().Format(time.RFC3339)
	fmt.Println(creationTimeStr)

	lat, _ := strconv.ParseFloat(eventcreationdata.Lat, 64)
	long, _ := strconv.ParseFloat(eventcreationdata.Long, 64)
	maxMem, _ := strconv.ParseInt(eventcreationdata.MaxMem, 10, 64)
	minMem, _ := strconv.ParseInt(eventcreationdata.MinMem, 10, 64)
	maxAge, _ := strconv.ParseInt(eventcreationdata.MaxAge, 10, 64)
	minAge, _ := strconv.ParseInt(eventcreationdata.MinAge, 10, 64)

	// Used for per user connection to DB
	dbconn := db.GetDBConn(DBName)
	defer dbconn.Close()

	// Add code to manage event creation request
	// Add an err handler here to ensure a failed signup request is handled
	stmt, _ := dbconn.Prepare(`INSERT INTO Events(event_name, lat, lng,
          creation_time, creator_id, start_time, end_time, max_mem, min_mem,
          friend_only, gender, min_age, max_age)
          VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);`)

	_, execerr := stmt.Exec(string(eventcreationdata.EventName),
		lat, long, creationTimeStr,
		string(eventcreationdata.Creatorid), eventcreationdata.StartTime,
		eventcreationdata.EndTime, maxMem, minMem, eventcreationdata.FriendOnly,
		eventcreationdata.Gender, minAge, maxAge)

	if execerr != nil {
		// If execution err occurs then throw error
		log.Fatal(execerr)
		ThrowInternalErrAndExit(w)
	}

	// If no error then give a success response
	RespondSuccessAndExit(w, "Event Created Successfully")

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

		event := generics.Event{

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

	// Printing the query to confirm output
	fmt.Println(query)

	dbconn := db.GetDBConn(DBName)
	res, err := dbconn.Query(query)
	if err != nil {
		panic(err)
	}

	return res, err
}
