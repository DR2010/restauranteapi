// Package main is the main package
// -------------------------------------
// .../restauranteapi/routes.go
// -------------------------------------
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route is
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes are
type Routes []Route

// VNewRouter is
func VNewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// XNewRouter is
func XNewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{"dishlist", "GET", "/dishlist", Hdishlist},
	Route{"dishfind", "GET", "/dishfind", Hdishfind},
	Route{"dishadd", "POST", "/dishadd", Hdishadd},
	Route{"orderadd", "POST", "/orderadd", Horderadd},
	Route{"dishupdate", "POST", "/dishupdate", Hdishupdate},
	Route{"dishdelete", "POST", "/dishdelete", Hdishdelete},
	Route{"dishalsolist", "GET", "/dishalsolist", Hdishalsolist},
	Route{"getcachedvalues", "GET", "/getcachedvalues", getcachedvalues},

	Route{"orderlist", "GET", "/orderlist", OrderList},
}
