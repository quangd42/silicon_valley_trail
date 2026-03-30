package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type ActionView struct {
	Kind model.Action
	Desc string
}

type PromptView struct {
	Actions []ActionView
}

func Prompt(cont *content.Content) PromptView {
	return PromptView{
		Actions: []ActionView{
			{Kind: model.ActionTravel, Desc: cont.Actions[model.ActionTravel].Desc},
			{Kind: model.ActionRest, Desc: cont.Actions[model.ActionRest].Desc},
			{Kind: model.ActionBuild, Desc: cont.Actions[model.ActionBuild].Desc},
			{Kind: model.ActionMarket, Desc: cont.Actions[model.ActionMarket].Desc},
			{Kind: model.ActionSave, Desc: cont.Actions[model.ActionSave].Desc},
			{Kind: model.ActionQuit, Desc: cont.Actions[model.ActionQuit].Desc},
		},
	}
}
