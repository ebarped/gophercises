package main

import (
	"fmt"
	"log"

	"phone/phone"
	"phone/store"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slices"
)

var numbers = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func main() {
	err := store.InitDB(numbers)
	if err != nil {
		log.Fatalf("error initializing the data store: %s\n", err)
	}

	fmt.Println("Normalizing numbers")
	err = normalizeNumbers()
	if err != nil {
		log.Fatalf("error normalizing numbers: %s\n", err)
	}

	fmt.Println("Removing duplicates")
	err = removeDuplicates()
	if err != nil {
		log.Fatalf("error removing duplicate numbers: %s\n", err)
	}
}

func normalizeNumbers() error {
	for _, p := range store.GetPhones() {
		normalizedNumber := phone.NormalizeNumber(p)
		fmt.Printf("updating phone number: %d. %q -> %q\n", p.Id, p.Number, normalizedNumber)
		err := store.UpdatePhoneNumber(p.Id, normalizedNumber.Number)
		if err != nil {
			return fmt.Errorf("error updating number %q to %q: %s\n", p.Number, normalizedNumber, err)
		}
	}
	return nil
}

func removeDuplicates() error {
	for _, p := range store.GetPhones() {
		currentNumbers := store.GetPhones() // get actual numbers, maybe we updated some of them in this loop
		currentNumbers = phone.Remove(currentNumbers, p)
		if slices.Contains(currentNumbers, p) {
			fmt.Printf("%s already exists in %+v, removing it (id=%d)\n", p.Number, currentNumbers, p.Id)
			err := store.RemovePhone(p.Id)
			if err != nil {
				return fmt.Errorf("error removing duplicate number %d. %q: %s\n", p.Id, p.Number, err)
			}
		}
	}
	return nil
}
