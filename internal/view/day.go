package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

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

type ActionResultView struct {
	Narative     []string
	Delta        model.Resources
	LocationName string
}

func ActionResult(ar logic.ActionResult, c *content.Content) ActionResultView {
	return ActionResultView{
		Narative:     c.Actions[ar.Action].Narrative,
		Delta:        ar.Delta,
		LocationName: ar.LocationName,
	}
}
