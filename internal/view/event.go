package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type EventView struct {
	Name      string
	Narrative []string
}

func Event(def gamedef.EventData) EventView {
	return EventView{
		Name:      def.Name,
		Narrative: def.Narrative,
	}
}

type EventResultView struct {
	Narrative []string
	Delta     model.Resources
}

func EventResult(eventIndex, choiceIndex int, er logic.EventResult, def gamedef.EventData) EventResultView {
	choiceDef := def.Choices[choiceIndex]
	return EventResultView{
		Narrative: choiceDef.Narrative,
		Delta:     er.Delta,
	}
}
