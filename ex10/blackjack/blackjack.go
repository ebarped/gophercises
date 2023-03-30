package blackjack

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ebarped/gophercises/ex9/deck"
)

type TurnOption int

const (
	HIT_OPTION TurnOption = iota + 1
	PASS_OPTION
)

type Game struct {
	deck    deck.Deck
	players []*Player
}

func (g Game) String() string {
	var result string
	result += "------------------------------\n"
	// print players status
	for _, p := range g.players {
		result += fmt.Sprintln(p)
	}

	// print deck status
	result += fmt.Sprintf("cards remaining in deck: %d\n", len(g.deck))

	result += "------------------------------"
	return result
}

// NewGame creates a new blackjack game with 1 shuffled deck and playerCount players (dealer included in this count)
func NewGame(playerCount int) Game {
	g := Game{
		deck: deck.New(deck.Shuffle),
	}

	for i := 0; i < playerCount-1; i++ {
		// create players except dealer
		g.players = append(g.players, &Player{
			name: fmt.Sprintf("player-%d", i+1),
		})

	}

	// create dealer
	g.players = append(g.players, &Player{
		name:     "dealer",
		isDealer: true,
	})

	return g
}

func (g Game) Start() {
	// initial deal for each player
	g.init()

	// main loop
	playerTurn := 0
	for {
		p := g.players[playerTurn]
		fmt.Println(g)

		fmt.Println("Turn of", p.name)

		var opt TurnOption
		if p.isDealer { // handle dealer logic
			p.cards[1].isVisible = true // turn up the face-down card
			opt = p.dealerIA()
		} else { // handle user logic
			fmt.Println("Options:")
			fmt.Println("1. Hit: get a card")
			fmt.Println("2. Stand: pass")
			var input int
			_, err := fmt.Scanf("%d", &input)
			if err != nil {
				fmt.Printf("error reading input from player: %s\n", err)
			}
			if !checkOption(input) {
				panic("option not allowed")
			}
			opt = TurnOption(input)
		}

		clearScreen()

		if opt == HIT_OPTION {
			c := g.deck.Deal()
			bc := NewBlackjackCard(c, true)
			p.cards = append(p.cards, bc)

			p.checkBusted()
			if p.isBusted {
				// dealer wins
				g.checkWinner()
			} else {
				// continue his turn
				continue
			}
		}

		// all players have played, check for winner
		playersCount := len(g.players[:]) - 1
		if playerTurn == playersCount {
			winner := g.checkWinner()
			if winner.name == "none" {
				fmt.Println("There is no winner, everyone is busted!")
				fmt.Println(g)
				os.Exit(0)
			}
			fmt.Println("Winner is", winner.name)
			fmt.Println(g)
			os.Exit(0)
		}

		playerTurn++
	}

}

// init executes the starting deals of the game
func (g *Game) init() {
	clearScreen()
	// deal first 2 cards to each player
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			c := g.deck.Deal()
			var visibility bool
			if p.isDealer && len(p.cards) == 1 { //give the dealer the last card face down
				visibility = false
			} else {
				visibility = true
			}
			bc := NewBlackjackCard(c, visibility)
			p.cards = append(p.cards, bc)
		}
	}
}

func (g Game) checkWinner() *Player {
	winner := Player{
		name: "none",
	}

	for _, p := range g.players {
		if !p.isBusted && p.score() > winner.score() {
			winner = *p
		}
	}
	return &winner
}

func (p Player) dealerIA() TurnOption {
	if p.score() <= 16 {
		return HIT_OPTION
	}
	return PASS_OPTION
}

func (p *Player) checkBusted() {
	if p.score() > 21 {
		p.isBusted = true
	}
}

func checkOption(opt int) bool {
	if opt != 1 && opt != 2 {
		return false
	}
	return true
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
