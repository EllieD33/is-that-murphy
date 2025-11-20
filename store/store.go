package store

import (
	"strings"
	"sync"

	"github.com/ellied33/is-that-murphy/models"
)


// Creating a look-up, allows multiple reads at once, but locks while writing happens
var (
	verifiedMap = map[string]models.VerifiedValue{}
	mu sync.RWMutex
)

// Lock while writing because writing alters memory, unsafe to allow multiple actions during this
func Add(v models.VerifiedValue) {
	mu.Lock()
	defer mu.Unlock()

	verifiedMap[strings.ToLower(v.Value)] = v
}

// Check whether entry exists, any number of reads can happen at once
func IsVerified(value string) (models.VerifiedValue, bool) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := verifiedMap[strings.ToLower(value)]
	return v, ok
}

// Reset helper for testing
func Reset() {
    mu.Lock()
    defer mu.Unlock()
    verifiedMap = map[string]models.VerifiedValue{}
}