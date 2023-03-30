package blackjack

import "fmt"

// this map translates the card number to its value in blackjack
var cardValues = map[int]int{
	1:  11, // ace, if the sum of the cards is >21, ace changes its value to 1
	2:  2,
	3:  3,
	4:  4,
	5:  5,
	6:  6,
	7:  7,
	8:  8,
	9:  9,
	10: 10,
	11: 10, // jack
	12: 10, // queen
	13: 10, // king
}

type Player struct {
	name     string
	cards    []BlackjackCard
	isDealer bool
	softHand bool // this hand has an ace, so the value of his hand can be "value" or "value-10"
	isBusted bool
}

func (p Player) String() string {
	var result string
	result += p.name
	result += " - "
	result += fmt.Sprint(p.cards)
	result += " - "
	result += fmt.Sprint("score: ", p.score())
	result += " - "
	result += fmt.Sprint("busted: ", p.isBusted)
	return result
}

// Score returns the score of the cards of the player
func (p *Player) score() int {
	var score int
	for _, c := range p.cards {
		if c.isVisible {
			score += cardValues[c.Value()]
		}
	}

	return score
}
