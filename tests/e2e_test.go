package test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"kiln-projects/api/routers"
	"kiln-projects/database"
	poller "kiln-projects/pollers"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// Mise en place d'un test complet avec un PostgreSQL de test qui est également déployé par le docker compose
// Il possède ses propres authentifiants et un schéma similaire à la DB de prod
type DelegationsResponse struct {
	Data []struct {
		Timestamp time.Time `json:"timestamp"`
		Amount    int64     `json:"amount"`
		Delegator string    `json:"delegator"`
		Level     int64     `json:"level"`
	} `json:"data"`
}

func TestEndToEndDelegations(t *testing.T) {
	t.Run("E2E Delegations", func(t *testing.T) {
		ctx := context.Background()
		// Chargement de mes variables d'environnement
		err := godotenv.Load("../.env")
		if err != nil {
			log.Println("no .env file found")
			return
		}

		// partie DB
		dbURL := os.Getenv("DB_URL_TEST")
		err = database.InitDB(ctx, dbURL)
		require.NoError(t, err)
		defer database.CloseDB()

		fakeTestDelegation := poller.Delegations{
			Timestamp: time.Date(2019, 6, 1, 12, 0, 0, 0, time.UTC),
			Amount:    123456789,
			Delegator: "tz1FakeAddress",
			Level:     42,
		}
		err = database.BulkAddingDelegations(ctx, []poller.Delegations{fakeTestDelegation})
		require.NoError(t, err)

		// Partie API
		router := routers.NewRouter()
		router.InitRoutes()
		server := httptest.NewServer(router)
		defer server.Close()

		// Test avec un query parameter sur l'année
		resp, err := http.Get(server.URL + "/xtz/delegations?year=2019")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// Partie récupération json et test contenu
		var result DelegationsResponse
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)
		require.NotEmpty(t, result.Data)

		yearTest := 2019
		for _, delegation := range result.Data {
			require.NotZero(t, delegation.Timestamp.Year())
			require.Equal(t, yearTest, delegation.Timestamp.Year())
			require.Greater(t, delegation.Amount, int64(0))
			require.Greater(t, delegation.Level, int64(0))
		}

	})
}
