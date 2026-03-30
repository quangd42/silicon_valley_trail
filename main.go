package main

import (
	"log"
	"os"

	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/program"
	"github.com/quangd42/silicon_valley_trail/internal/save"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
)

func main() {
	renderer := ui.NewTerminal(os.Stdin, os.Stdout)
	saver := save.NewJSONSaver("") // TODO: load config
	gameCopy := content.Load()
	state := model.NewState(content.DefaultRoute())
	err := program.Run(
		renderer,
		saver,
		gameCopy,
		state,
	)
	if err != nil {
		log.Fatal("internal error")
	}
}
