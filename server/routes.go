package server

import "net/http"

func (s *Server) SetupRoutes() {
	s.Router.Use(setContentTypeFileFormat)
	s.Router.HandleFunc("/", s.serveIndex)
	s.Router.HandleFunc("/{blog:blog\\/?}", s.serveMarkdownIndex)
	s.Router.HandleFunc("/lab", s.serveMarkdownIndex)
	s.Router.HandleFunc("/systems", s.serveMarkdownIndex)
	s.Router.HandleFunc("/blog/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/lab/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/systems/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/notepad", s.serveNotepad)
	s.Router.HandleFunc("/login", s.serveLogin)
	s.Router.HandleFunc("/todo", s.middlewareTokenVerification(s.serveTodo))
	s.Router.HandleFunc("/api/user/login/", setContentTypeJSON(s.handleLogin)).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/todo/", setContentTypeJSON(s.middlewareTokenVerification(s.handleTodoCreate))).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/todo/", setContentTypeJSON(s.middlewareTokenVerification(s.handleTodoUpdate))).Methods(http.MethodPut)
	s.Router.HandleFunc("/api/todo/", setContentTypeJSON(s.middlewareTokenVerification(s.handleTodoDelete))).Methods(http.MethodDelete)
	s.Router.HandleFunc("/api/todo/all", setContentTypeJSON(s.middlewareTokenVerification(s.handleTodoPending))).Methods(http.MethodGet)
	s.Router.HandleFunc("/notes", s.middlewareTokenVerification(s.handleNotesIndex)).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/note", setContentTypeJSON(s.middlewareTokenVerification(s.handleNoteContent))).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/note", setContentTypeJSON(s.middlewareTokenVerification(s.handleNoteCreate))).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/note", setContentTypeJSON(s.middlewareTokenVerification(s.handleNoteUpdate))).Methods(http.MethodPut)
	s.Router.HandleFunc("/health", s.handleHealth)
	s.Router.NotFoundHandler = http.HandlerFunc(s.serve404)
	s.setupStatic("static")
}
