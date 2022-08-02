package server

func (s *Server) SetupRoutes() {
	s.Router.HandleFunc("/", s.serveIndex)
	s.setupStatic("static")
}
