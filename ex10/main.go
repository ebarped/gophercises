package main

import (
	"ex10/blackjack"
)

const PLAYER_COUNT = 2 // 1 dealer & 1 player

func main() {
	g := blackjack.NewGame(PLAYER_COUNT)
	g.Start()
}
