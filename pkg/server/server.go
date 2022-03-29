package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port    string
	handler http.Handler
}

func (s *Server) Run() {
	host := fmt.Sprintf(":%s", s.port)

	httpServer := &http.Server{
		Addr:    host,
		Handler: s.handler,
	}

	log.Printf("[Start] Server started at %s", host)
	httpServer.ListenAndServe()
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{port: port, handler: handler}
}
