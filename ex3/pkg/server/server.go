package server

import (
	"ex3/pkg/history"
	"net/http"
	"text/template"
)

type server struct {
	router *http.ServeMux
	tpl    *template.Template
	h      history.History
}

func New(templatePath string, h history.History) *server {
	s := &server{
		router: http.NewServeMux(),
		tpl:    template.Must(template.ParseFiles(templatePath)),
		h:      h,
	}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
