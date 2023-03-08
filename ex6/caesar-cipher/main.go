package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
 * Complete the 'caesarCipher' function below.
 *
 * The function is expected to return a STRING.
 * The function accepts following parameters:
 *  1. STRING s
 *  2. INTEGER k
 */
const (
	ALPHABET_LEN = 26
	BASE_CHAR    = 97
)

func caesarCipher(s string, k int32) string {
	var result string

	for _, char := range s {

		if !unicode.IsLetter(char) {
			result += string(char)
			continue
		}

		cipheredChar := (char-BASE_CHAR+k)%ALPHABET_LEN + BASE_CHAR
		result += string(cipheredChar)
	}

	return result
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	fmt.Printf("input string s: ")
	s := readLine(reader)

	fmt.Printf("input rotation k: ")
	kTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	k := int32(kTemp)

	result := caesarCipher(s, k)

	fmt.Fprintf(os.Stdout, "result: %s\n", result)
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
