// Package gamedef defines the authored game data and scripted effects.
package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type Definition struct {
	Intro       Narrative
	Route       []model.Location
	Actions     map[model.Action]ActionData
	ActionOrder []model.Action
	Weather     map[model.WeatherKind]WeatherData
	Events      map[string]EventData
	EventPools  model.EventPools
	Endings     map[logic.Ending]EndingCopy
}

type Narrative []string

// Load returns the authored game definition.
func Load() *Definition {
	return &Definition{
		Intro:       introCopy(),
		Route:       DefaultRoute(),
		Actions:     actionData(),
		ActionOrder: actionOrder(),
		Weather:     weatherCopy(),
		Events:      eventData(),
		EventPools:  makeEventPools(AllEvents),
		Endings:     endingCopy(),
	}
}
