package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	cp "github.com/otiai10/copy"
)

const (
	SAMPLE_DIR = "sample"
	DIR        = "test"
)

func main() {
	// objective: change birthday_001.txt -> Birthday - 1 of 4.txt

	// copy "sample" dir into a new one to let the original files untouched
	err := cp.Copy(SAMPLE_DIR, DIR)
	if err != nil {
		panic(err)
	}

	// traverse all the files
	err = filepath.WalkDir(DIR, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Printf("Analyzing %q\n", path)
			if match(path) {
				newName := rename(path)
				fmt.Printf("Changing %q to %q\n", path, newName)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// only analyze strings that match
// (name)_(number).(extension)
func match(src string) bool {
	if strings.Contains(src, "_") {

		return true
	}
	return false
}

func rename(src string) string {
	// --- origin ---
	// birthday_001.txt
	// (name)_(number).(extension)
	// --- destination ---
	// Birthday - 1 of 4.txt
	// (name_upper) - (n1) of (n2).(extension)

	// obtain the filename without the path
	withoutPath := strings.SplitN(src, "/", 99)
	fileName := withoutPath[len(withoutPath)-1]

	// obtain the name, the number and the extension
	name := strings.Split(fileName, "_")[0]
	name = strings.Title(name)
	number := strings.Split(strings.Split(fileName, "_")[1], ".")[0]
	num, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
	extension := strings.Split(fileName, ".")[1]

	fmt.Printf("\tname=%q, number=%d, extension=%q\n", name, num, extension)
	return fmt.Sprintf("%s - %d of 4.%s", name, num, extension)
}
