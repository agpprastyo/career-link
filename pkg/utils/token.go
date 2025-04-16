package utils

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Use a package-level random source that's properly seeded once
var (
	tokenRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randMutex sync.Mutex // Protect the random generator for concurrent use
)

// Generate6DigitToken generates a random 6-digit token
func Generate6DigitToken() string {
	randMutex.Lock()
	defer randMutex.Unlock()
	return fmt.Sprintf("%06d", tokenRand.Intn(1000000))
}

// GenerateUniqueToken generates a token and ensures it doesn't exist in the database
func GenerateUniqueToken(ctx context.Context, repo func(ctx context.Context, token string) (bool, error)) (string, error) {
	for i := 0; i < 5; i++ { // Try up to 5 times
		token := Generate6DigitToken()
		exists, err := repo(ctx, token)
		if err != nil {
			return "", err
		}
		if !exists {
			return token, nil
		}
	}
	return "", errors.New("failed to generate unique token")
}
