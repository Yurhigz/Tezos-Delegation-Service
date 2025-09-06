package handlers

import (
	"kiln-projects/database"
	poller "kiln-projects/pollers"
	"net/http"
)

func GetDelegations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// vérifier qu'on a que les deux paramètres timestamp et level/blockheight
	params := map[string]bool{
		"timestamp":   true,
		"blockheight": true,
	}
	for key := range query {
		if !params[key] {
			http.Error(w, "invalid query parameters, only available parameters are timestamp and level", http.StatusBadRequest)
			return
		}
	}

	timestamp := query.Get("timestamp")
	if timestamp == "" {

	}
	blockheight := query.Get("blockheight")
	var DelegationsList []poller.Delegations

	database.DelegationsRetrieval(r.Context(), timestamp, blockheight)

}
