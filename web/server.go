package web

import "net/http"

type Service struct {
	r *http.ServeMux
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
