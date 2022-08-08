package server

func (s *Server) SetupRoutes() {
	s.Router.HandleFunc("/", s.serveIndex)
	s.Router.HandleFunc("/blog", s.serveMarkdownIndex)
	s.Router.HandleFunc("/lab", s.serveMarkdownIndex)
	s.Router.HandleFunc("/systems", s.serveMarkdownIndex)
	s.Router.HandleFunc("/blog/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/lab/{title}", s.serveMarkdownPost)
	s.Router.HandleFunc("/systems/{title}", s.serveMarkdownPost)
	s.setupStatic("static")
}
