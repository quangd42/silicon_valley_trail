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
	RenderEventInfo(view.EventView)
	RenderEventResult(v view.EventResultView)
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
		turn := p.playTurn(state)
		if turn.quitToMenu {
			return
		}
		if p.evaluateEnding(state) {
			return
		}
		if turn.traveled {
			// Event has a chance to play out after arriving at a new location
			// (beginning of a new day/turn).
			if p.playEvent(state) {
				return
			}
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

// Returns true when `QuitToMainMenu` is selected
func (p *Program) playEvent(state *model.State) bool {
	if state.CurrentLocation == 0 {
		return false
	}
	// Events have a 1 in 3 chance of occuring
	roll := p.rng.IntN(3)
	if roll >= 1 {
		return false
	}
	events := p.def.Events
	eventIndex := p.rng.IntN(len(events))
	eventDef := events[eventIndex]
	for {
		p.renderer.RenderEventInfo(view.Event(eventDef))
		choice := p.renderer.RenderPrompt(view.EventChoicePrompt(eventDef))
		switch choice.Kind {
		case model.ChoiceControl:
			if p.applyControl(state, choice.Control) {
				return true
			}
		case model.ChoiceEvent:
			result := p.applyEventChoice(state, choice.EventChoiceIndex, eventDef)
			p.renderer.RenderEventResult(view.EventResult(eventIndex, choice.EventChoiceIndex, result, eventDef))
			return false
		default:
			panic("non event choice from event")

		}
	}
}

type TurnInfo struct {
	quitToMenu bool
	traveled   bool
}

func (p *Program) playTurn(state *model.State) TurnInfo {
	p.refreshWeather(state)
	for {
		p.renderer.RenderDayInfo(view.Day(state, p.def))
		choice := p.renderer.RenderPrompt(view.DayPrompt(p.def))
		switch choice.Kind {
		case model.ChoiceControl:
			if p.applyControl(state, choice.Control) {
				return TurnInfo{quitToMenu: true}
			}
		case model.ChoiceAction:
			res := p.applyAction(state, choice.Action)
			p.renderer.RenderActionResult(view.ActionResult(choice.Action, res, p.def))
			return TurnInfo{traveled: choice.Action == model.ActionTravel}
		default:
			panic("non action choice from day action")
		}
	}
}

// Returns true when game losing conditions are met
func (p *Program) evaluateEnding(state *model.State) bool {
	ending := logic.EvaluateEnding(state)
	if ending != logic.EndingNone {
		p.renderer.RenderEnding(view.Ending(ending, p.def))
		return true
	}
	return false
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

func (p *Program) applyAction(state *model.State, action model.Action) logic.ActionResult {
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

func (p *Program) applyEventChoice(state *model.State, choiceIndex int, def gamedef.EventData) logic.EventResult {
	choice := def.Choices[choiceIndex]
	return logic.ApplyEventChoiceEffect(state, choice.Effect)
}

func (p *Program) applyControl(state *model.State, control model.Control) bool {
	switch control {
	case model.ControlSave:
		p.saveGame(state)
		return false
	case model.ControlQuitToMenu:
		p.renderer.ClearScreen()
		return true
	default:
		panic("invalid in-game session control")
	}
}
