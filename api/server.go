package api

import (
	"github.com/gin-gonic/gin"

	db "tech-school/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	return &Server{
		store:  store,
		router: gin.Default(),
	}
}

func (s *Server) Run(addr string) error {
	s.initRoutes()

	return s.router.Run(addr)
}

func (s *Server) initRoutes() {
	s.router.GET("/accounts", s.listAccounts)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.POST("/accounts", s.createAccount)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
