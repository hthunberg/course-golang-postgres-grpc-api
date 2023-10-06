package random

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/pkg/currency"
)

// Int generates a random integer between 0 and max.
// It will panic if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func Int(max int64) int64 {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}

	return randomNumber.Int64()
}

// String returns a random generated string of length n.
// It will panic if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func String(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			panic(err)
		}
		ret[i] = alphabet[num.Int64()]
	}

	return string(ret)
}

// Owner generates a random owner name.
// It will panic if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func Owner() string {
	return String(6)
}

// Money generates a random amount of money between 0 and 1000.
// It will panic if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func Money() int64 {
	return Int(1000)
}

// Currency generates a random currency code.
// It will panic if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func Currency() string {
	currencies := []string{currency.USD, currency.EUR, currency.SEK}
	n := int64(len(currencies))
	return currencies[Int(n)]
}

// Email generates a random email.
// It will panic if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func Email() string {
	return fmt.Sprintf("%s@email.com", String(6))
}
