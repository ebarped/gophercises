package server

import "log"

func (s *server) routes() {

	log.Printf("templateHandler: registering %q handler\n", "/")
	s.router.HandleFunc("/", s.TemplateHandler())
}
