package controllers

import (
	"go-lbapp/api"
	"net/http"

	"github.com/gorilla/mux"
)

var restroutes = Routes{
	Route{
		"user",
		"GET",
		"/v1/user/{userid}",
		api.UserHandler,
	},
}

// SetResourceRoutes : Used to add REST resource routes to mux
func SetResourceRoutes(router *mux.Router) *mux.Router {
	for _, route := range restroutes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = APILogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
