package model

import "net/http"

// KeyPair : A struct holding a pair of key and value
type KeyPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Route : Generic format for all routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes : generic slice for all routes
type Routes []Route
