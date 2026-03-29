package view

import "github.com/quangd42/silicon_valley_trail/internal/model"

type DayView struct {
	Day       int
	Resources model.Resources
	Party     model.Party
	Location  model.Location
	Progress  int
	Weather   model.WeatherKind
}

func Day(s *model.State) DayView {
	return DayView{
		Day:       s.Day,
		Resources: s.Resources,
		Party:     s.Party,
		Location:  s.Route[s.CurrentLocation],
		Progress:  (s.CurrentLocation + 1) * 100 / len(s.Route),
		Weather:   s.Weather,
	}
}

func StandardActions() [6]model.ActionKind {
	return [6]model.ActionKind{
		model.ActionTravel,
		model.ActionRest,
		model.ActionBuild,
		model.ActionMarket,
		model.ActionSave,
		model.ActionQuit,
	}
}
