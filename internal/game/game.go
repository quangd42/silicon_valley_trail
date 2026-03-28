package game

import (
	"fmt"
	"io"
	"os"
)

type Game struct {
	state    string
	renderer io.Writer
	store    map[string]string
}

func NewGame() *Game {
	return &Game{
		state:    "this is a brand new game",
		renderer: os.Stdout,
		store:    nil,
	}
}

func (g *Game) Run() error {
	fmt.Fprintf(g.renderer, "current state: %s\n", g.state)
	return nil
}
