package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ellied33/is-that-murphy/handlers"
	"github.com/ellied33/is-that-murphy/models"
	"github.com/ellied33/is-that-murphy/store"
)

func resetStore() { store.Reset() }

func TestVerifyHandler(t *testing.T) {
	resetStore()

	v := models.VerifiedValue{
		Value: "test@test.com",
		Type: "email",
	}

	store.Add(v)

	req, err := http.NewRequest("GET", "/verify?value=test@test.com", nil)

	if err != nil {
    	t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.VerifyHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var got models.VerifiedValue
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if got.Value != v.Value || got.Type != v.Type {
    	t.Errorf("unexpected result: got %+v want %+v", got, v)
	}
}

func TestVerifyHandler_TrimsWhitespace(t *testing.T) {
    resetStore()

    v := models.VerifiedValue{
        Value: "test@test.com",
        Type:  "email",
    }
    store.Add(v)

    req := httptest.NewRequest("GET", "/verify", nil)
	q := req.URL.Query()
	q.Set("value", "   test@test.com   ")
	req.URL.RawQuery = q.Encode()

    rr := httptest.NewRecorder()

    handlers.VerifyHandler(rr, req)
    res := rr.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Fatalf("expected status 200, got %d", res.StatusCode)
    }

    var got models.VerifiedValue
    if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
        t.Fatalf("failed to decode response JSON: %v", err)
    }

    if got.Value != v.Value || got.Type != v.Type {
        t.Errorf("unexpected result: got %+v want %+v", got, v)
    }
}

func TestVerifyHandler_NormalisesInput(t *testing.T) {
    resetStore()

    v := models.VerifiedValue{
        Value: "test@test.com",
        Type:  "email",
    }
    store.Add(v)

    req := httptest.NewRequest("GET", "/verify?value=test@Test.Com", nil)

    rr := httptest.NewRecorder()

    handlers.VerifyHandler(rr, req)
    res := rr.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Fatalf("expected status 200, got %d", res.StatusCode)
    }

    var got models.VerifiedValue
    if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
        t.Fatalf("failed to decode response JSON: %v", err)
    }

    if got.Value != v.Value || got.Type != v.Type {
        t.Errorf("unexpected result: got %+v want %+v", got, v)
    }
}

func TestVerifyHandler_NotFound(t *testing.T) {
	resetStore()

	req := httptest.NewRequest("GET", "/verify?value=unknown", nil)
	rr := httptest.NewRecorder()

	handlers.VerifyHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if res["type"] != "not verified" {
		t.Errorf("expected type 'not verified', got %v", res["type"])
	}
}

func TestVerifyHandler_MissingValue(t *testing.T) {
	resetStore()

	req := httptest.NewRequest("GET", "/verify", nil)
	rr := httptest.NewRecorder()

	handlers.VerifyHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rr.Code)
	}
}

func TestAddHandler_Success(t *testing.T) {
	resetStore()

	payload := models.VerifiedValue{
		Value: "doggo@murphy.com",
		Type:  "email",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handlers.AddHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	var got models.VerifiedValue
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if got != payload {
		t.Errorf("expected %+v but got %+v", payload, got)
	}

	if _, ok := store.IsVerified("doggo@murphy.com"); !ok {
		t.Fatalf("store should contain value after POST")
	}
}

func TestAddHandler_InvalidJSON(t *testing.T) {
	resetStore()

	req := httptest.NewRequest("POST", "/verify", strings.NewReader("{not-json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handlers.AddHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rr.Code)
	}
}
