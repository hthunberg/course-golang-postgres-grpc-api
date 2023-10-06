//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSomeIntegration(t *testing.T) {
	err := testDbInstance.Ping(context.Background())
	require.NoError(t, err)
}
