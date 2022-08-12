package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/viveknathani/nattukaka/service"
)

type Server struct {
	Router  *mux.Router
	Service *service.Service
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
