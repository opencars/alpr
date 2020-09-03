package http

func (s *server) configureRouter() {
	router := s.router.PathPrefix("/api/v1/alpr").Subrouter()

	router.Handle("/public/version", s.Version()).Methods("GET", "OPTIONS")
	router.Handle("/private/recognize", s.Recognize()).Methods("GET", "OPTIONS")
}
