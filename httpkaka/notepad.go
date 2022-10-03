package httpkaka

import (
	"io/ioutil"
	"net/http"
)

func (s *Server) serveNotepad(w http.ResponseWriter, r *http.Request) {

	html, err := ioutil.ReadFile("static/pages/notepad.html")
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}
