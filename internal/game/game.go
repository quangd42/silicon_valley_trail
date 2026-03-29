package game

import (
	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
	"github.com/quangd42/silicon_valley_trail/internal/view"
)

func Run(
	copy *content.Content,
	r *ui.Terminal,
	state *model.State,
) error {
	standardActions := view.StandardActions()
	r.RenderIntro(view.IntroView(copy.Intro))
	for state.CurrentLocation < len(state.Route)-1 {
		r.RenderDay(view.Day(state))
		action := r.Prompt(standardActions[:])
		msg := logic.ApplyAction(state, action)
		r.RenderInfo(msg)
	}
	r.RenderEnding(view.EndingView(copy.Ending))
	return nil
}
