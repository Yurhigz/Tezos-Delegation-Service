package main

import (
	"context"
	"fmt"
	"kiln-projects/api/routers"
	"kiln-projects/database"
	poller "kiln-projects/pollers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func StartPoller(ctx context.Context, baseURL string) {
	url := baseURL
	for {
		delegationBatch, err := poller.PollTzkt(url)
		if err != nil {
			log.Printf("ERR | polling data: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		if len(delegationBatch) == 0 {
			log.Println("No new delegations, waiting…")
			//Eventuellement possibilité de faire une progression dans l'attente avec un intervalle max d'attente
			time.Sleep(10 * time.Minute)
			continue
		}

		err = database.BulkAddingDelegations(ctx, delegationBatch)
		if err != nil {
			log.Printf("ERR | inserting into DB: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		offset := delegationBatch[len(delegationBatch)-1].Level
		url = fmt.Sprintf("%s&level.gt=%d", baseURL, offset)

		log.Printf("Inserted %d delegations, next offset=%d", len(delegationBatch), offset)

		time.Sleep(5 * time.Second)
	}
}

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
		return
	}
	dbURL := os.Getenv("DB_URl")

	if err := database.InitDB(ctx, dbURL); err != nil {
		log.Printf("ERR   | %v", err)
		return
	}

	baseUrl := os.Getenv("TZKT_API_ENDPOINT")

	go StartPoller(ctx, baseUrl)

	router := routers.NewRouter()
	router.InitRoutes()

	err = http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatalf("ERR | server: %v", err)
	}

}
