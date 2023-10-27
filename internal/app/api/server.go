package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/bank"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	bank   bank.Bank
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, bank bank.Bank) (*Server, error) {
	server := &Server{
		config: config,
		bank:   bank,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/accounts", server.createAccount)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse formats the errors returned to the client.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
