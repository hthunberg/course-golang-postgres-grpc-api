package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/bank"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, bank bank.Bank) *Server {
	config := util.Config{}

	server, err := NewServer(config, bank)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
