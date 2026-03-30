package program

import (
	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
	"github.com/quangd42/silicon_valley_trail/internal/view"
)

func Run(
	cont *content.Content,
	r *ui.Terminal,
	state *model.State,
) error {
	r.RenderIntro(view.IntroView(cont.Intro))

	for state.CurrentLocation < len(state.Route)-1 {
		r.RenderDay(view.Day(state))
		selection := r.PromptSelection(view.InGamePrompt(cont))
		if selection.Kind {
			res := logic.ApplyAction(state, selection.Action)
			r.RenderActionResult(view.ActionResult(res, cont))
		} else {
			switch selection.Control {
			case model.ControlSave:
				saveGame(r, state)
			case model.ControlQuitToMenu:
				quitToMenu(r)
			default:
				panic("invalid in-game control")
			}
		}
	}
	r.RenderEnding(view.EndingView(cont.Ending))
	return nil
}

func saveGame(r *ui.Terminal, state *model.State) {
	// Save state
	_ = state
	// Render info
	r.RenderInfo("Game saved.\n")
}

func quitToMenu(r *ui.Terminal) {
	// Render Intro menu
	// PromptSelection
	r.RenderInfo("Menu here.\n")
}
