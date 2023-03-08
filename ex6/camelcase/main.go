package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

/*
 * Complete the 'camelcase' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts STRING s as parameter.
 */

func camelcase(s string) int32 {
	// marks if a new word should start (avoids inputs with 2 consecutive
	// uppercase letters)
	newWord := false

	wordCount := 0

	if unicode.IsUpper(rune(s[0])) {
		fmt.Fprintln(os.Stderr, "error: the first letter has to be lowercase...")
		os.Exit(1)
	}

	wordCount++

	for i, char := range s {

		// check that we only work with letters
		if !unicode.IsLetter(char) {
			fmt.Fprintln(os.Stderr, "error: you have entered something that is not a letter...")
			os.Exit(1)
		}

		// we should be starting with lowercase now, last word is finished
		if newWord {
			if unicode.IsUpper(char) {
				fmt.Fprintln(os.Stderr, "error: it should not be 2 consecutive uppercase letters...")
				fmt.Printf("%c at index %d\n", char, i)
				os.Exit(1)
			}
			newWord = false
		}

		if unicode.IsUpper(char) {
			wordCount++
			newWord = true
		}

	}

	return int32(wordCount)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	fmt.Printf("Introduce a camelcase string: ")
	s := readLine(reader)

	result := camelcase(s)

	fmt.Println("Words:", result)
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
