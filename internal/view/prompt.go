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
	Actions       []ActionView
	ActionsLabel  string
	Controls      []model.Control
	ControlsLabel string
}

func InGamePrompt(cont *content.Content) PromptView {
	return PromptView{
		Actions: []ActionView{
			{Kind: model.ActionTravel, Desc: cont.Actions[model.ActionTravel].Desc},
			{Kind: model.ActionRest, Desc: cont.Actions[model.ActionRest].Desc},
			{Kind: model.ActionBuild, Desc: cont.Actions[model.ActionBuild].Desc},
			{Kind: model.ActionMarket, Desc: cont.Actions[model.ActionMarket].Desc},
		},
		ActionsLabel: "Actions:\n",
		Controls: []model.Control{
			model.ControlSave,
			model.ControlQuitToMenu,
		},
		ControlsLabel: "Controls:\n",
	}
}
