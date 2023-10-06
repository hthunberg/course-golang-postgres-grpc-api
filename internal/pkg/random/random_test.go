//go:build !integration

package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	// A stupid test, just me elaborating with crypto/rand
	// The case is that i wanted to assure that we do have 3
	// entries in map ["SEK", "USD", "EUR"].
	currencies := make(map[string]string, 10)
	for i := 0; i < 100; i++ {
		currencies[Currency()] = "whatever"
	}
	assert.Len(t, currencies, 3)
}
