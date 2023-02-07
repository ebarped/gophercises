package main

import (
	"ex4/pkg/link"
	"fmt"
	"io"
	"log"
	"os"
)

const fileName = "ex4.html"

func main() {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error opening file %s: %s\n", fileName, err)
	}
	defer file.Close()

	links, err := link.GetLinks(file)
	if err != nil && err != io.EOF {
		log.Fatalf("error getting Link Tokens from file %s: %s\n", file.Name(), err)
	}

	for _, link := range links {
		fmt.Println(link)
	}

}
