package server

import (
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) serveNotepad(w http.ResponseWriter, r *http.Request) {

	html, err := ioutil.ReadFile("static/pages/notepad.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}
