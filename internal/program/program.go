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
		action := r.Prompt(view.Prompt(cont))
		res := logic.ApplyAction(state, cont, action)
		r.RenderInfo(res.Msg)
	}
	r.RenderEnding(view.EndingView(cont.Ending))
	return nil
}
