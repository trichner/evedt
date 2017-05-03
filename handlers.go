package evedt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/trichner/evedt/tracker"
)

func Index(repl *tracker.Replicator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome!")
	}
}

func DonationsIndex(repl *tracker.Replicator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		limit := 100
		days := 30

		limits := params["limit"]
		if len(limits) == 1 {
			limit, _ = strconv.Atoi(limits[0])
		}

		daysStr := params["days"]
		if len(daysStr) == 1 {
			days, _ = strconv.Atoi(daysStr[0])
		}

		donations := repl.FindDonations(limit, days)

		appendJson(w, donations)
	}
}

func DonationsTop(repl *tracker.Replicator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		limit := 100
		days := 30

		limits := params["limit"]
		if len(limits) == 1 {
			limit, _ = strconv.Atoi(limits[0])
		}

		daysStr := params["days"]
		if len(daysStr) == 1 {
			days, _ = strconv.Atoi(daysStr[0])
		}
		fmt.Printf("Days: %d\n", days)
		donations := repl.FindTopDonations(limit, days)

		appendJson(w, donations)
	}
}
func appendJson(w http.ResponseWriter, r interface{}) error {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(r)
}
