package poller

import (
	"encoding/json"
	"net/http"
	"time"
)

// Block Height = level
type Delegations struct {
	Timestamp time.Time `json:"timestamp"`
	Amount    int64     `json:"amount"`
	Delegator string    `json:"delegator"`
	Level     int64     `json:"level"`
}
type RawDelegations struct {
	Timestamp time.Time    `json:"timestamp"`
	Amount    int64        `json:"amount"`
	Sender    RawDelegator `json:"sender"`
	Level     int64        `json:"level"`
}

type RawDelegator struct {
	Alias   string `json:"alias,omitempty"`
	Address string `json:"address"`
}

func PollTzkt(url string) ([]Delegations, error) {
	var rawDelegationsList []RawDelegations
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// log.Println("got a response")

	err = json.NewDecoder(resp.Body).Decode(&rawDelegationsList)
	if err != nil {
		return nil, err
	}

	DelegationList := make([]Delegations, 0, len(rawDelegationsList))
	for _, d := range rawDelegationsList {
		DelegationList = append(DelegationList, Delegations{
			Timestamp: d.Timestamp,
			Amount:    d.Amount,
			Delegator: d.Sender.Address,
			Level:     d.Level,
		})
	}
	return DelegationList, nil
}
