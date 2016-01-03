package controllers

import (
	"go-lbapp/api"
	"net/http"

	"github.com/gorilla/mux"
)

var restroutes = Routes{
	Route{
		"getUser",
		"GET",
		"/v1/user/{userid}",
		api.UserHandler,
	},
	Route{
		"getEvent",
		"GET",
		"/v1/event/{eventid}",
		api.GetEventHandler,
	},
	Route{
		"setEvent",
		"POST",
		"/v1/event",
		api.CreateEvent,
	},
	Route{
		"deleteEvent",
		"DELETE",
		"/v1/event/{eventid}",
		api.DeleteEvent,
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
