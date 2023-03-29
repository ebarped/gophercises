package deck

import (
	"fmt"
	"math/rand"
	"strconv"

	"golang.org/x/exp/slices"
)

type Suit int

const (
	SuitNone Suit = iota // jokers has no suit
	SuitDiamonds
	SuitClubs
	SuitHearts
	SuitSpades
)

const (
	DiamondsRepr = '♦'
	ClubsRepr    = '♣'
	HeartsRepr   = '♥'
	SpadesRepr   = '♠'
	JokerRepr    = '🃏'
)

type Card struct {
	value int
	suit  Suit
}

func NewCard(val int, suit Suit) Card {
	return Card{
		value: val,
		suit:  suit,
	}
}

func (c Card) Value() int {
	return c.value
}

func (c Card) String() string {
	var SuitRepr rune
	var valueRepr string

	switch c.suit {
	case SuitNone:
		SuitRepr = JokerRepr
	case SuitDiamonds:
		SuitRepr = DiamondsRepr
	case SuitClubs:
		SuitRepr = ClubsRepr
	case SuitHearts:
		SuitRepr = HeartsRepr
	case SuitSpades:
		SuitRepr = SpadesRepr
	}

	switch c.value {
	case 0:
		valueRepr = "Jk"
	case 1:
		valueRepr = "A"
	case 11:
		valueRepr = "J"
	case 12:
		valueRepr = "Q"
	case 13:
		valueRepr = "K"
	default:
		valueRepr = strconv.Itoa(c.value)
	}

	return fmt.Sprintf("(%s,%c)", valueRepr, SuitRepr)
}

type Deck []Card

func (d Deck) String() string {
	var result string
	for _, card := range d {
		result += fmt.Sprintln(card)
	}
	return result
}

// New returns, by default, an ordered deck of cards
// you can pass options to modify the result
func New(opts ...func(Deck) Deck) Deck {
	var deck Deck
	suits := []Suit{SuitDiamonds, SuitClubs, SuitHearts, SuitSpades}

	for _, suit := range suits {
		for j := 0; j < 13; j++ {
			deck = append(deck, Card{
				value: j + 1,
				suit:  suit,
			})
		}
	}

	// apply opts
	for _, opt := range opts {
		deck = opt(deck)
	}

	return deck
}

// Shuffle shuffles the deck
func Shuffle(d Deck) Deck {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
	return d
}

// WithJokers adds jokers to the deck, honoring the signature of the opts func needed by the New constructor
func WithJokers(n int) func(Deck) Deck {
	return func(d Deck) Deck {
		for i := 0; i < n; i++ {
			d = append(d, Card{
				value: 0,
				suit:  SuitNone,
			})
		}
		return d
	}
}

// WithFilteredCards removes from the deck all the cards with the value specified
func WithFilteredCards(values ...int) func(Deck) Deck {
	return func(d Deck) Deck {
		for _, val := range values {
			for i, c := range d {
				if c.value == val {
					d = slices.Delete(d, i, i+1)
				}
			}
		}
		return d
	}

}

// WithAdditionalDecks allows to add additional decks to the original
func WithAdditionalDecks(count int) func(Deck) Deck {
	return func(d Deck) Deck {
		for i := 0; i < count; i++ {
			d = append(d, New()...)
		}
		return d
	}
}

// RemoveCard removes card by its index from the deck
func (d *Deck) removeCard(i int) {
	*d = slices.Delete(*d, i, i+1)

}

// Deal deals a count cards and removes them from the deck
func (d *Deck) Deal() Card {
	c := (*d)[0]
	d.removeCard(0)

	return c
}
