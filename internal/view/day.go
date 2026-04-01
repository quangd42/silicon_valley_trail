package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type DayView struct {
	Day           int
	Resources     model.Resources
	Party         model.Party
	Location      model.Location
	Progress      int
	Weather       model.WeatherKind
	WeatherImpact string
}

func Day(s *model.State, def *gamedef.Definition) DayView {
	return DayView{
		Day:           s.Day,
		Resources:     s.Resources,
		Party:         s.Party,
		Location:      s.Route[s.CurrentLocation],
		Progress:      (s.CurrentLocation + 1) * 100 / len(s.Route),
		Weather:       s.Weather,
		WeatherImpact: def.Weather[s.Weather].Desc,
	}
}

type ActionResultView struct {
	Narative     []string
	Delta        model.Resources
	LocationName string
	Weather      model.WeatherKind
}

func ActionResult(ar logic.Result, def *gamedef.Definition) ActionResultView {
	// CurrentLocation will always >= 1 after a Travel action, so we're reusing
	// default value 0 as sentinel value for "did not travel".
	var locName string
	if ar.CurrentLocation == 0 {
		locName = ""
	} else {
		locName = def.Route[ar.CurrentLocation].Name
	}
	return ActionResultView{
		Narative:     def.Actions[ar.Action].Narrative,
		Delta:        ar.Delta,
		LocationName: locName,
		Weather:      ar.Weather,
	}
}
