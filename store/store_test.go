package store_test

import (
	"testing"

	"github.com/ellied33/is-that-murphy/models"
	"github.com/ellied33/is-that-murphy/store"
)


func TestAddAndVerify(t *testing.T) {
	store.Reset()

	v := models.VerifiedValue{
		Value: "test@test.com",
		Type: "email",
	}

	store.Add(v)

	got, ok := store.IsVerified("test@test.com")

	if !ok {
		t.Fatalf("expected value to exist but got ok = false")
	}

	if got.Value != v.Value {
		t.Errorf("expected value %q but got %q", v.Value, got.Value)
	}

	if got.Type != v.Type {
		t.Errorf("expected type %q but got %q", v.Type, got.Type)
	}
}

func TestLookupIsCaseInsensitive(t *testing.T) {
	store.Reset()

	v := models.VerifiedValue{
		Value: "DOGGO123",
		Type:  "twitter",
	}

	store.Add(v)

	_, ok := store.IsVerified("dOgGo123")
	if !ok {
		t.Fatalf("expected case-insensitive match but lookup failed")
	}
}

func TestLookupMissingValue(t *testing.T) {
	store.Reset()

	_, ok := store.IsVerified("not-here")
	if ok {
		t.Fatalf("expected lookup to fail but got ok = true")
	}
}