package server

import "net/http"

func (s *Server) SetupRoutes() {
	s.Router.Use(setContentTypeFileFormat)
	s.Router.HandleFunc("/", s.serveIndex)
	s.Router.HandleFunc("/blog", s.serveMarkdownIndex)
	s.Router.HandleFunc("/lab", s.serveMarkdownIndex)
	s.Router.HandleFunc("/systems", s.serveMarkdownIndex)
	s.Router.HandleFunc("/blog/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/lab/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/systems/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/notepad", s.serveNotepad)
	s.Router.HandleFunc("/login", s.serveLogin)
	s.Router.HandleFunc("/api/user/login/", setContentTypeJSON(s.handleLogin)).Methods(http.MethodPost)
	s.setupStatic("static")
}
