package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"kiln-projects/database"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
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

func BulkAddingDelegations(parentsContext context.Context, DelegationsList []Delegations) error {
	ctx, cancel := context.WithTimeout(parentsContext, 10*time.Second)
	defer cancel()

	_, err := database.DBPool.CopyFrom(ctx, pgx.Identifier{"delegations"}, []string{"Timestamp", "SenderAddress", "Amount", "BlockHeight"}, pgx.CopyFromSlice(len(DelegationsList), func(i int) ([]any, error) {
		return []any{DelegationsList[i].Timestamp, DelegationsList[i].Sender.Address, DelegationsList[i].Amount, DelegationsList[i].BlockHeight}, nil
	}))

	if err != nil {
		return fmt.Errorf("ERR | Error inserting delegations : %v", err)
	}

	log.Printf("%d Delegations added successfully, last BlockHeight %v", len(DelegationsList), DelegationsList[len(DelegationsList)-1].BlockHeight)
	return nil

}
