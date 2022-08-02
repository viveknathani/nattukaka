package server

import (
	"html/template"
	"log"
	"net/http"
)

func (s *Server) setupStatic(directory string) {

	fileServer := http.FileServer(http.Dir(directory))
	s.Router.PathPrefix("/" + directory + "/").Handler(http.StripPrefix("/"+directory, fileServer))
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	p, err := template.ParseFiles("static/pages/home.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = p.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
