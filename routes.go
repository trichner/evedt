package evedt

import (
	"net/http"
	"github.com/trichner/evedt/tracker"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	NewHandler func (r *tracker.Replicator) http.HandlerFunc
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
