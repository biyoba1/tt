package mdk

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    200 * time.Second,
		WriteTimeout:   200 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}
