package httpkaka

import (
	"html/template"
	"net/http"
)

func (s *Server) setupStatic(directory string) {

	fileServer := http.FileServer(http.Dir(directory))
	s.Router.PathPrefix("/" + directory + "/").Handler(http.StripPrefix("/"+directory, fileServer))
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	p, err := template.ParseFiles("static/pages/home.html")
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}
	err = p.Execute(w, nil)
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}
}
