package evedt

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"DonationsIndex",
		"GET",
		"/donations",
		DonationsIndex,
	},
	Route{
		"DonationsTop",
		"GET",
		"/donations/top",
		DonationsTop,
	},
}
