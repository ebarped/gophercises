package main

import (
	"flag"
	"fmt"
)

func main() {
	csvFilePath := "problems.csv"
	csvFile := flag.String("csvFile", csvFilePath, "path to csv file")
	flag.Parse()

	fmt.Println("Loading questions from " + *csvFile)
	return

}
