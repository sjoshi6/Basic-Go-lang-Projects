package api

import "net/http"

// Route : Generic format for all routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes : generic slice for all routes
type Routes []Route
