package program

import (
	"context"
	"fmt"
	"os"

	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/view"
)

type Renderer interface {
	RenderIntro(view.IntroView)
	RenderDayInfo(view.DayView)
	RenderPrompt(view.PromptView) model.PromptChoice
	RenderActionResult(view.ActionResultView)
	RenderInfo(string)
	RenderInfoNoWait(string)
	RenderEnding(view.EndingView)
	ClearScreen()
}

type Saver interface {
	Save(*model.State) error
	Load(*model.State) error
}

type WeatherProvider interface {
	WeatherAt(context.Context, model.Location) (model.WeatherKind, error)
}

type Program struct {
	renderer Renderer
	saver    Saver
	weather  WeatherProvider
	rng      logic.RNG
	def      *gamedef.Definition
	exit     func(int)
}

func New(
	renderer Renderer,
	saver Saver,
	weather WeatherProvider,
	rng logic.RNG,
	def *gamedef.Definition,
) *Program {
	return &Program{
		renderer: renderer,
		saver:    saver,
		weather:  weather,
		rng:      rng,
		def:      def,
		exit:     os.Exit,
	}
}

func (p *Program) Run() {
	for {
		selection := p.renderer.RenderPrompt(view.MainMenuPrompt())
		if selection.Kind != model.ChoiceControl {
			panic("in-game action on main menu")
		}
		switch selection.Control {
		case model.ControlNewGame:
			p.newGame()
		case model.ControlLoad:
			p.loadGame()
		case model.ControlQuitGame:
			p.quitGame()
		default:
			panic("invalid game session control on main menu")
		}
	}
}

func (p *Program) startGame(state *model.State, isNew bool) {
	p.renderer.ClearScreen()
	if isNew {
		p.renderer.RenderIntro(view.IntroView(p.def.Intro))
	}
	for state.CurrentLocation < len(state.Route)-1 {
		p.refreshWeather(state)
		p.renderer.RenderDayInfo(view.Day(state, p.def))
		selection := p.renderer.RenderPrompt(view.DayPrompt(p.def))
		if selection.Kind == model.ChoiceAction {
			res := p.applyAction(state, selection.Action)
			p.renderer.RenderActionResult(view.ActionResult(res, p.def))
		} else {
			switch selection.Control {
			case model.ControlSave:
				p.saveGame(state)
			case model.ControlQuitToMenu:
				// simply return from the game loop because we're
				// already in the main menu loop
				p.renderer.ClearScreen()
				return
			default:
				panic("invalid game session control")
			}
		}
		ending := logic.EvaluateEnding(state)
		if ending != logic.EndingNone {
			p.renderer.RenderEnding(view.Ending(ending, p.def))
			return
		}
	}

	ending := logic.ResolveFinalEnding(state, p.rng)
	p.renderer.RenderEnding(view.Ending(ending, p.def))
}

func (p *Program) newGame() {
	state := model.NewState(gamedef.DefaultRoute())
	p.startGame(state, true)
}

func (p *Program) loadGame() error {
	var state model.State
	err := p.saver.Load(&state)
	if err != nil {
		p.renderer.RenderInfo(fmt.Sprintf("Failed to load game: %s", err.Error()))
		return err
	}
	p.startGame(&state, false)
	return nil
}

func (p *Program) saveGame(state *model.State) {
	err := p.saver.Save(state)
	if err != nil {
		p.renderer.RenderInfo(fmt.Sprintf("Failed to save game: %s", err.Error()))
		return
	}
	p.renderer.RenderInfo("Game saved.")
}

func (p *Program) quitGame() {
	p.renderer.RenderInfoNoWait("Bye!")
	exit := p.exit
	if exit == nil {
		exit = os.Exit
	}
	exit(0)
}

func (p *Program) refreshWeather(state *model.State) {
	if p.weather == nil || len(state.Route) == 0 {
		state.Weather = model.WeatherUnknown
		return
	}

	weather, err := p.weather.WeatherAt(context.Background(), state.Route[state.CurrentLocation])
	if err != nil {
		state.Weather = model.WeatherUnknown
		return
	}
	state.Weather = weather
}

func (p *Program) applyAction(state *model.State, action model.Action) logic.Result {
	actionDef := p.def.Actions[action]
	weatherDef := p.def.Weather[state.Weather]
	return logic.ApplyActionEffects(
		state,
		action,
		actionDef.Effect,
		weatherDef.Effect,
		p.rng,
	)
}
