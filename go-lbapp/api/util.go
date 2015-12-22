package api

import (
	"encoding/json"
	"net/http"
)

/*  All util functions for API Calls
    Typically used to send JSON replies back to the client
    Fixed format calls for 200,400 & 500 status codes
*/

// RespondOrThrowErr : Respond to general requests or exit with server err.
func RespondOrThrowErr(responseObj BasicResponse, w http.ResponseWriter) {

	responseJSON, err := json.Marshal(responseObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// ThrowInternalErrAndExit : Respond with internal server error
func ThrowInternalErrAndExit(w http.ResponseWriter) {

	responsecontent := BasicResponse{
		"Internal Server Error",
		500,
	}

	w.Header().Set("StatusCode", "500")
	w.Header().Set("Status", "Internal Server Error")
	RespondOrThrowErr(responsecontent, w)
}

// RespondSuccessAndExit : Repond with a success
func RespondSuccessAndExit(w http.ResponseWriter, msg string) {

	responsecontent := BasicResponse{
		msg,
		200,
	}
	w.Header().Set("StatusCode", "200")
	RespondOrThrowErr(responsecontent, w)

}
