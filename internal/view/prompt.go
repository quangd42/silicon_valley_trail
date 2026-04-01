package view

import (
	"fmt"

	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type PromptItemView struct {
	Text   string
	Choice model.PromptChoice
}

type PromptSectionView struct {
	Label string
	Items []PromptItemView
}

type PromptView struct {
	Title    string
	Sections []PromptSectionView
}

func DayPrompt(def *gamedef.Definition) PromptView {
	return PromptView{Sections: []PromptSectionView{
		{
			Label: "Actions:",
			Items: actionItems(def, def.ActionOrder),
		},
		{
			Label: "Controls:",
			Items: controlItems([]model.Control{
				model.ControlSave,
				model.ControlQuitToMenu,
			}),
		},
	}}
}

func EventChoicePrompt(event gamedef.EventData) PromptView {
	return PromptView{Sections: []PromptSectionView{
		{
			Label: event.ChoicesLabel,
			Items: eventChoiceItems(event.Choices),
		},
		{
			Label: "Controls:",
			Items: controlItems([]model.Control{
				model.ControlSave,
				model.ControlQuitToMenu,
			}),
		},
	}}
}

func actionItems(def *gamedef.Definition, actions []model.Action) []PromptItemView {
	items := make([]PromptItemView, 0, len(actions))
	for _, action := range actions {
		items = append(items, PromptItemView{
			Text: def.Actions[action].Desc,
			Choice: model.PromptChoice{
				Kind:   model.ChoiceAction,
				Action: action,
			},
		})
	}
	return items
}

func controlItems(controls []model.Control) []PromptItemView {
	items := make([]PromptItemView, 0, len(controls))
	for _, control := range controls {
		items = append(items, PromptItemView{
			Text: control.String(),
			Choice: model.PromptChoice{
				Kind:    model.ChoiceControl,
				Control: control,
			},
		})
	}
	return items
}

func eventChoiceItems(choices []gamedef.EventChoiceData) []PromptItemView {
	items := make([]PromptItemView, 0, len(choices))
	for i, choice := range choices {
		items = append(items, PromptItemView{
			Text: eventChoiceText(choice),
			Choice: model.PromptChoice{
				Kind:             model.ChoiceEvent,
				EventChoiceIndex: i,
			},
		})
	}
	return items
}

func eventChoiceText(choice gamedef.EventChoiceData) string {
	switch {
	case choice.Name == "":
		return choice.Desc
	case choice.Desc == "" || choice.Desc == "...":
		return choice.Name
	default:
		return fmt.Sprintf("%s (%s)", choice.Name, choice.Desc)
	}
}
