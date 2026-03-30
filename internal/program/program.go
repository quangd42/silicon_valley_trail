package program

import (
	"fmt"
	"os"

	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/save"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
	"github.com/quangd42/silicon_valley_trail/internal/view"
)

func Run(
	renderer *ui.Terminal,
	saver save.Saver,
	cont *content.Content,
) {
	for {
		selection := renderer.RenderMainMenu(view.MainMenu())
		if selection.Kind {
			panic("in-game action on main menu")
		}
		switch selection.Control {
		case model.ControlNewGame:
			newGame(renderer, saver, cont)
		case model.ControlLoad:
			loadGame(renderer, saver, cont)
		case model.ControlQuitGame:
			quitGame(renderer)
		default:
			panic("invalid game session control on main menu")
		}
	}
}

func startGame(
	rndr *ui.Terminal,
	saver save.Saver,
	cont *content.Content,
	state *model.State,
	new bool,
) {
	rndr.ClearScreen()
	if new {
		rndr.RenderIntro(view.IntroView(cont.Intro))
		rndr.ClearScreen()
	}
	for state.CurrentLocation < len(state.Route)-1 {
		rndr.RenderDay(view.Day(state))
		selection := rndr.PromptSelection(view.InGamePrompt(cont))
		if selection.Kind {
			res := logic.ApplyAction(state, selection.Action)
			rndr.RenderActionResult(view.ActionResult(res, cont))
		} else {
			switch selection.Control {
			case model.ControlSave:
				saveGame(rndr, saver, state)
			case model.ControlQuitToMenu:
				// simply return from the game loop because we're
				// already in the main menu loop
				rndr.ClearScreen()
				return
			default:
				panic("invalid game session control")
			}
		}
		ending := logic.EvaluateEnding(state)
		if ending != logic.EndingNone {
			rndr.RenderEnding(view.Ending(ending, cont))
			return
		}
	}

	ending := logic.ResolveFinalEnding(state)
	rndr.RenderEnding(view.Ending(ending, cont))
}

func newGame(
	rndr *ui.Terminal,
	saver save.Saver,
	cont *content.Content,
) {
	state := model.NewState(content.DefaultRoute())
	startGame(rndr, saver, cont, state, true)
}

func loadGame(
	rndr *ui.Terminal,
	saver save.Saver,
	cont *content.Content,
) error {
	var state model.State
	err := saver.Load(&state)
	if err != nil {
		rndr.RenderInfo(fmt.Sprintf("Failed to load game: %s", err.Error()))
		return err
	}
	startGame(rndr, saver, cont, &state, false)
	return nil
}

func saveGame(
	rndr *ui.Terminal,
	saver save.Saver,
	state *model.State,
) {
	err := saver.Save(state)
	if err != nil {
		rndr.RenderInfo(fmt.Sprintf("Failed to save game: %s", err.Error()))
		return
	}
	rndr.RenderInfo("Game saved.")
}

func quitGame(rndr *ui.Terminal) {
	rndr.RenderInfoNoWait("Bye!")
	os.Exit(0)
}
