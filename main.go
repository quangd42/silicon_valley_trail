package main

import (
	"log"

	"github.com/quangd42/silicon_valley_trail/internal/game"
)

func main() {
	game := game.NewGame()
	if err := game.Run(); err != nil {
		log.Fatal("internal error")
	}
}
