package server

import (
	"fmt"
	"net/http"
)

func (s server) TemplateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		if len(queryParams) == 0 {
			s.tpl.Execute(w, s.h["intro"])
			return
		}

		chapter := queryParams.Get("chapter")
		if val, ok := s.h[chapter]; ok {
			fmt.Printf("vamos al cap %s\n", val)
			s.tpl.Execute(w, s.h[chapter])
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("The chapter %q does not exists!", chapter)))
	})
}
