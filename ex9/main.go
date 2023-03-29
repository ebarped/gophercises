package main

import (
	"fmt"

	"github.com/ebarped/gophercises/ex9/deck"
)

func main() {
	fmt.Println("----------- DEFAULT DECK ------------------")
	d1 := deck.New()
	fmt.Println(d1)

	fmt.Println("------------ SHUFFLED DECK -----------------")
	d3 := deck.New(deck.Shuffle)
	fmt.Println(d3)

	fmt.Println("------------ DECK WITH JOKERS -----------------")
	d4 := deck.New(deck.WithJokers(4))
	fmt.Println(d4)

	fmt.Println("------------ DECK WITHOUT 2s & 3s -----------------")
	d5 := deck.New(deck.WithFilteredCards(2, 3))
	fmt.Println(d5)

	fmt.Println("------------ DECK WITH 2 DECKS -----------------")
	d6 := deck.New(deck.WithAdditionalDecks(1))
	fmt.Println(d6)

	fmt.Println("------------ DEAL CARDS -----------------")
	d7 := deck.New()
	c1 := d7.Deal()
	c2 := d7.Deal()
	fmt.Println("card dealt:", c1)
	fmt.Println("card dealt:", c2)
	fmt.Println(d7)

}
