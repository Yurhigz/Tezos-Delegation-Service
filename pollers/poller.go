package poller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Block Height = level

type Delegations struct {
	Timestamp   time.Time `json:"timestamp"`
	Sender      Sender    `json:"sender"`
	Amount      int64     `json:"amount"`
	BlockHeight int64     `json:"level"`
}

type Sender struct {
	Alias   string `json:"alias,omitempty"`
	Address string `json:"address"`
}

func PollTzkt(url string) ([]Delegations, error) {
	var DelegationsList []Delegations
	resp, err := http.Get(url)
	if err != nil {
		return DelegationsList, err
	}
	defer resp.Body.Close()
	log.Println("got a response")

	err = json.NewDecoder(resp.Body).Decode(&DelegationsList)
	if err != nil {
		return DelegationsList, err
	}

	return DelegationsList, nil
}
