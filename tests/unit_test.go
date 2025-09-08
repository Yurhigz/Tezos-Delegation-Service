package test

import (
	"fmt"
	poller "kiln-projects/pollers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Tester 3 cas :
// - JSON valide
// - JSON invalide
// - API injoignable ou une erreur réseau
func TestPollTzkt(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []poller.Delegations
		expectError bool
	}{
		{
			name:  "Valid JSON",
			input: `[{"timestamp":"2025-09-08T12:00:00Z","amount":1000,"sender":{"alias":"Alice","address":"tz1..."},"level":100}]`,
			expected: []poller.Delegations{
				{
					Timestamp: time.Date(2025, 9, 8, 12, 0, 0, 0, time.UTC),
					Amount:    1000,
					Delegator: "tz1...",
					Level:     100,
				},
			},
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			input:       "INVALID_JSON",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "API unreachable",
			input:       "",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tt.input)
			}))
			defer ts.Close()

			result, err := poller.PollTzkt(ts.URL)

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.expectError {
				if len(result) != len(tt.expected) {
					t.Fatalf("expected %d items, got %d", len(tt.expected), len(result))
				}
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("expected %+v, got %+v", tt.expected[i], result[i])
					}
				}
			}
		})
	}
}
