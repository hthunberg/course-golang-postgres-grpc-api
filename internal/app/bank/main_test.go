//go:build !integration

package bank

import (
	"log"
	"os"
	"testing"
)

// TestMain allows us to run arbitrary code before and after tests run.
func TestMain(m *testing.M) {
	log.Println("BEFORE the tests")
	exitVal := m.Run()
	log.Println("AFTER the tests")
	os.Exit(exitVal)
}
