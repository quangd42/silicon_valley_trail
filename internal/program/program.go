package program

import (
	"context"
	"fmt"
	"os"

	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/save"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
	"github.com/quangd42/silicon_valley_trail/internal/view"
	"github.com/quangd42/silicon_valley_trail/internal/weather"
)

func Run(
	renderer *ui.Terminal,
	saver save.Saver,
	weather weather.Service,
	cont *content.Content,
) {
	for {
		selection := renderer.RenderMainMenu(view.MainMenu())
		if selection.Kind {
			panic("in-game action on main menu")
		}
		switch selection.Control {
		case model.ControlNewGame:
			newGame(renderer, saver, weather, cont)
		case model.ControlLoad:
			loadGame(renderer, saver, weather, cont)
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
	weather weather.Service,
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
		refreshWeather(state, weather)
		rndr.RenderDay(view.Day(state, cont))
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
	weather weather.Service,
	cont *content.Content,
) {
	state := model.NewState(content.DefaultRoute())
	startGame(rndr, saver, weather, cont, state, true)
}

func loadGame(
	rndr *ui.Terminal,
	saver save.Saver,
	weather weather.Service,
	cont *content.Content,
) error {
	var state model.State
	err := saver.Load(&state)
	if err != nil {
		rndr.RenderInfo(fmt.Sprintf("Failed to load game: %s", err.Error()))
		return err
	}
	startGame(rndr, saver, weather, cont, &state, false)
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

func refreshWeather(state *model.State, svc weather.Service) {
	if svc == nil || len(state.Route) == 0 {
		state.Weather = model.WeatherUnknown
		return
	}

	current, err := svc.Current(context.Background(), state.Route[state.CurrentLocation])
	if err != nil {
		state.Weather = model.WeatherUnknown
		return
	}
	state.Weather = current
}
