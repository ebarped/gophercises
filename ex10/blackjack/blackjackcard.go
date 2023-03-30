package blackjack

import "github.com/ebarped/gophercises/ex9/deck"

// wrapper of deck.Card to add visibility
type BlackjackCard struct {
	deck.Card
	isVisible bool
}

func NewBlackjackCard(c deck.Card, v bool) BlackjackCard {
	return BlackjackCard{
		Card:      c,
		isVisible: v,
	}
}

func (bc BlackjackCard) String() string {
	if bc.isVisible {
		return bc.Card.String()
	}
	return "(X,X)"
}
