package main

import (
	"ex3/pkg/history"
	"ex3/pkg/server"
	"log"
	"net/http"
	"strconv"
)

const (
	port        = 8080
	tplFile     = "templates/web.html"
	historyFile = "gopher.json"
)

func main() {

	history, err := history.New(historyFile)
	if err != nil {
		panic(err)
	}

	s := server.New(tplFile, history)
	log.Printf("Starting the server on :%d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), s))

}
