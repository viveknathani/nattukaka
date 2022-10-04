package httpkaka

import "net/http"

func (s *Server) serve404(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<h3>oops, 404!</h3>"))
}
