package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellied33/is-that-murphy/handlers"
	"github.com/ellied33/is-that-murphy/models"
	"github.com/ellied33/is-that-murphy/store"
)

func resetStore() { store.Reset() }

func TestVerifyHandler(t *testing.T) {
	cases := []struct {
		name       string
		setup      func()
		queryValue string
		wantStatus int
		wantValue  string
		wantType   string
	}{
		{
			name: "Found exact match",
			setup: func() {
				resetStore()
				store.Add(models.VerifiedValue{Value: "test@test.com", Type: "email"})
			},
			queryValue: "test@test.com",
			wantStatus: http.StatusOK,
			wantValue:  "test@test.com",
			wantType:   "email",
		},
		{
			name: "Trims whitespace",
			setup: func() {
				resetStore()
				store.Add(models.VerifiedValue{Value: "test@test.com", Type: "email"})
			},
			queryValue: "   test@test.com   ",
			wantStatus: http.StatusOK,
			wantValue:  "test@test.com",
			wantType:   "email",
		},
		{
			name: "Not found",
			setup: func() {
				resetStore()
			},
			queryValue: "unknown@test.com",
			wantStatus: http.StatusOK,
			wantValue:  "unknown@test.com",
			wantType:   "not verified",
		},
		{
			name: "Missing value",
			setup: func() {
				resetStore()
			},
			queryValue: "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			req := httptest.NewRequest("GET", "/verify", nil)

			if tc.queryValue != "" {
				q := req.URL.Query()
				q.Set("value", tc.queryValue)
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.VerifyHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d", tc.wantStatus, rr.Code)
			}

			if tc.wantStatus == http.StatusOK {
				var got models.VerifiedValue
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode JSON: %v", err)
				}

				if got.Value != tc.wantValue || got.Type != tc.wantType {
					t.Errorf("unexpected result: got %+v want value=%q type=%q", got, tc.wantValue, tc.wantType)
				}
			}
		})
	}
}

func TestAddHandler(t *testing.T) {
	cases := []struct {
		name       string
		payload    any
		wantStatus int
		wantStore  bool
	}{
		{
			name: "Successful add",
			payload: models.VerifiedValue{
				Value: "doggo@murphy.com",
				Type:  "email",
			},
			wantStatus: http.StatusCreated,
			wantStore:  true,
		},
		{
			name:       "Invalid JSON",
			payload:    "{not-json}",
			wantStatus: http.StatusBadRequest,
			wantStore:  false,
		},
		{
			name:       "Empty JSON object",
			payload:    models.VerifiedValue{},
			wantStatus: http.StatusCreated,
			wantStore:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetStore()

			var bodyBytes []byte
			switch v := tc.payload.(type) {
			case string:
				bodyBytes = []byte(v)
			case models.VerifiedValue:
				b, err := json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal payload: %v", err)
				}
				bodyBytes = b
			default:
				t.Fatalf("unsupported payload type")
			}

			req := httptest.NewRequest("POST", "/verify", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.AddHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d", tc.wantStatus, rr.Code)
			}

			if tc.wantStatus == http.StatusCreated {
				var got models.VerifiedValue
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response JSON: %v", err)
				}

				if _, ok := store.IsVerified(got.Value); !ok && tc.wantStore {
					t.Fatalf("store should contain value after POST")
				}
			}
		})
	}
}
