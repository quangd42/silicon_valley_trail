package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
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

func InGamePrompt(def *gamedef.Definition) PromptView {
	return PromptView{
		Actions: []ActionView{
			{Kind: model.ActionTravel, Desc: def.Actions[model.ActionTravel].Desc},
			{Kind: model.ActionRest, Desc: def.Actions[model.ActionRest].Desc},
			{Kind: model.ActionBuild, Desc: def.Actions[model.ActionBuild].Desc},
			{Kind: model.ActionMarket, Desc: def.Actions[model.ActionMarket].Desc},
		},
		ActionsLabel: "Actions:\n",
		Controls: []model.Control{
			model.ControlSave,
			model.ControlQuitToMenu,
		},
		ControlsLabel: "Controls:\n",
	}
}
