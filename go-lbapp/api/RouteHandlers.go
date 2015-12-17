package api

import (
	"encoding/json"
	"net/http"
)

/* Contains all the Route Handlers for API function calls */

// SignUp : Handle Signup requests for new user
func SignUp(w http.ResponseWriter, r *http.Request) {

	responsecontent := BasicResponse{
		"User Registered Successfully",
		200,
	}

	respondOrThrowErr(responsecontent, w)
}

// Login : Handle Login requests for existing users
func Login(w http.ResponseWriter, r *http.Request) {

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
