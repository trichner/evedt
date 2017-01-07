package evedt

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func DonationsIndex(w http.ResponseWriter, r *http.Request) {

	limit := 20
	days := 20
	donations := repo.FindDonations(limit, days)

	appendJson(w, donations)
	w.WriteHeader(http.StatusOK)
}

func DonationsTop(w http.ResponseWriter, r *http.Request) {

	limit := 20
	days := 20
	donations := repo.FindTopDonations(limit, days)

	appendJson(w, donations)
	w.WriteHeader(http.StatusOK)
}

func appendJson(w http.ResponseWriter, r interface{}) error {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(r)
}

// func (api *API) abort(rw http.ResponseWriter, statusCode int) {
// 	rw.WriteHeader(rw)
// }
