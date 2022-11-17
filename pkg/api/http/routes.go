package http

func (s *server) configureRouter() {
	router := s.router.PathPrefix("/api/v1/alpr").Subrouter()

	router.Handle("/private/recognize", s.Recognize()).Methods("GET", "OPTIONS")
}
