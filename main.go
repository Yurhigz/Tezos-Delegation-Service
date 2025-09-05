package main

import (
	"context"
	"kiln-projects/database"
	poller "kiln-projects/pollers"
	"log"
	"time"
)

func main() {

	dbURL := "host=127.0.0.1 user=postgres password=postgres dbname=tzktdb port=5442 sslmode=disable TimeZone=Asia/Shanghai"
	parentCtx := context.Background()
	err := database.InitDB(parentCtx, dbURL)

	if err != nil {
		log.Printf("ERR   | %v", err)
		return
	}

	var CurrentDelegationsBatch []poller.Delegations
	var Offset int
	url := "https://cryptoslam.api.tzkt.io/v1/operations/delegations?select=timestamp,sender,amount,level&limit=200"
	go func() {
		for {
			CurrentDelegationsBatch, err := poller.PollTzkt(url)
			if err != nil {
				log.Printf("ERR | Error polling data : %w", err)
				time.Sleep(15 * time.Second)
				continue
			}
			// stockage de CurrentDelegationsBatch dans la DB puis remise à zéro de la value et update du pointeur de sauvegarde
			err = poller.BulkAddingDelegations(parentCtx, CurrentDelegationsBatch)
			if err != nil {
				log.Printf("ERR | Error bulk adding data in the DB : %w", err)
				time.Sleep(15 * time.Second)
				continue
			}
			Offset = int(CurrentDelegationsBatch[len(CurrentDelegationsBatch)-1].BlockHeight)
			time.Sleep(time.Minute)
		}

	}()

}
