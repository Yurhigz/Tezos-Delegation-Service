package test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kiln-projects/api/routers"
	"kiln-projects/database"
	poller "kiln-projects/pollers"

	"github.com/stretchr/testify/require"
)

//

func TestEndToEndDelegations(t *testing.T) {
	ctx := context.Background()

	// Init DB (tu peux pointer vers une DB de test en docker)
	dbURL := "host=localhost user=postgres password=postgres dbname=tzktdb_test port=5442 sslmode=disable"
	err := database.InitDB(ctx, dbURL)
	require.NoError(t, err)

	// Insert fake data
	fakeDelegation := poller.Delegations{
		Timestamp: time.Now(),
		Amount:    123456789,
		Delegator: "tz1FakeAddress",
		Level:     42,
	}
	err = database.BulkAddingDelegations(ctx, []poller.Delegations{fakeDelegation})
	require.NoError(t, err)

	// Start API in test mode
	router := routers.NewRouter()
	router.InitRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	// Call API
	resp, err := http.Get(server.URL + "/xtz/delegations?level=42")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string][]poller.Delegations
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	require.Len(t, result["Delegations"], 1)
	require.Equal(t, fakeDelegation.Delegator, result["Delegations"][0].Delegator)
}
