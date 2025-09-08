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
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func StartPoller(ctx context.Context, baseURL string) {
	url := baseURL

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping poller")
			return
		default:
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
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
		return
	}
	dbURL := os.Getenv("DB_URL")

	if err := database.InitDB(ctx, dbURL); err != nil {
		log.Printf("ERR   | %v", err)
		return
	}
	defer database.CloseDB()

	baseUrl := os.Getenv("TZKT_API_ENDPOINT")

	go StartPoller(ctx, baseUrl)

	router := routers.NewRouter()
	router.InitRoutes()
	srv := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERR | server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Termination signal received, shutting down server...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	cancel() // Pour mettre fin au poller
}
