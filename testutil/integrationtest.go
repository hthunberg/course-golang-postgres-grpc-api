package testutil

import (
	"os"
	"testing"
)

// Investigating the possibilities to not use build tags for separating different
// type of tests, see https://konradreiche.com/blog/how-to-separate-integration-tests-in-go

func IntegrationTest(t *testing.T) {
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
}
