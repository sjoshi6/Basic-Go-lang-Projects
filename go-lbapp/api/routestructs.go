package api

import "net/http"

// Route : Common struct for all routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes : Generic Array Struct for all routes
type Routes []Route

// BasicResponse : JSON reply for API Calls
type BasicResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
