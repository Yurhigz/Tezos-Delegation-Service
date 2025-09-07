package handlers

import (
	"encoding/json"
	"kiln-projects/database"
	poller "kiln-projects/pollers"
	"net/http"
	"strconv"
	"time"
)

func GetDelegations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// vérifier qu'on a que les deux paramètres timestamp et level/blockheight
	params := map[string]bool{
		"timestamp": true,
		"level":     true,
	}
	for key := range query {
		if !params[key] {
			http.Error(w, "invalid query parameters, only available parameters are timestamp and level", http.StatusBadRequest)
			return
		}
	}
	var timestampValue int

	timestamp := query.Get("timestamp")
	if timestamp != "" {
		parsedTime, err := time.Parse("2006", timestamp)
		if err != nil {
			http.Error(w, "invalid timestamp format: must be a year (YYYY)", http.StatusBadRequest)
			return
		}
		timestampValue = parsedTime.Year()
	} else {
		timestampValue = time.Now().Year()
	}

	level := query.Get("level")
	var levelValue int64
	if level != "" {
		val, err := strconv.Atoi(level)
		if err != nil {
			http.Error(w, "invalid number as level : must be an int superior to 0", http.StatusBadRequest)
			return
		}
		levelValue = int64(val)
	}

	DelegationsList, err := database.DelegationsRetrieval(r.Context(), timestampValue, levelValue)

	if err != nil {
		http.Error(w, "delegations retrieval error reaching the DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]poller.Delegations{
		"Data": DelegationsList,
	})
}
