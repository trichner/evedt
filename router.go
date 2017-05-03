package evedt

import (
	"github.com/gorilla/mux"
	"github.com/trichner/evedt/tracker"
)

// Creates a new router under the specified prefix
// and registers all routes
func NewRouter(prefix string, r *tracker.Replicator) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	if len(prefix) > 0 {
		router = router.PathPrefix(prefix).Subrouter()
	}

	for _, route := range routes {
		// decorate with logger
		handler := Logger(route.NewHandler(r), route.Name)

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
