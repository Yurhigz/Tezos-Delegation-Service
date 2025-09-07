package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test pour PollTzkt - le plus important à tester car il fait l'appel HTTP
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

	result, err := PollTzkt(server.URL)

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

	result, err := PollTzkt(server.URL)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test pour GetDelegations - cas nominal
func TestGetDelegations_ValidRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/xtz/delegations", nil)
	w := httptest.NewRecorder()

	// Note: ce test nécessiterait un mock de la DB ou une DB de test
	// Pour l'exercice, on peut le laisser en commentaire ou utiliser une DB en mémoire

	// GetDelegations(w, req)
	// assert.Equal(t, 200, w.Code)
	// assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	t.Skip("Test nécessite une DB de test - à implémenter selon l'infrastructure choisie")
}

// Test pour GetDelegations - paramètres invalides
func TestGetDelegations_InvalidParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/xtz/delegations?invalid_param=test", nil)
	w := httptest.NewRecorder()

	GetDelegations(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid query parameters")
}

// Test pour GetDelegations - année invalide
func TestGetDelegations_InvalidYear(t *testing.T) {
	req := httptest.NewRequest("GET", "/xtz/delegations?timestamp=invalid", nil)
	w := httptest.NewRecorder()

	GetDelegations(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid timestamp format")
}

// Test pour GetDelegations - level invalide
func TestGetDelegations_InvalidLevel(t *testing.T) {
	req := httptest.NewRequest("GET", "/xtz/delegations?level=not-a-number", nil)
	w := httptest.NewRecorder()

	GetDelegations(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid number as level")
}
