package controllers

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
