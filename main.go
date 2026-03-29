package main

import (
	"log"
	"os"

	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/game"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
)

func main() {
	gameCopy := content.Load()
	renderer := ui.NewTerminal(os.Stdin, os.Stdout)
	state := model.NewState(content.DefaultRoute())
	err := game.Run(gameCopy, renderer, state)
	if err != nil {
		log.Fatal("internal error")
	}
}
