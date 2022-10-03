package httpkaka

import "net/http"

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		return
	}
}
