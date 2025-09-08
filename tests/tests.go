package test

import (
	poller "kiln-projects/pollers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPollTzkt(t *testing.T) {
	// Test avec une réponse valide
	mockResponse := `[
		{
			"timestamp": "2024-01-01T12:00:00Z",
			"amount": 1000,
			"level": 100,
			"sender": {
				"address": "tz1abc123"
			}
		}
	]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	result, err := poller.PollTzkt(server.URL)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "tz1abc123", result[0].Delegator)
	assert.Equal(t, int64(1000), result[0].Amount)
	assert.Equal(t, int64(100), result[0].Level)
}

// Test pour PollTzkt avec JSON invalide
func TestPollTzkt_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	result, err := poller.PollTzkt(server.URL)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test du BulkAdding dans la DB

// Test du router/handler

// Test du retrieval
