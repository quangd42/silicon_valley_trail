package logic

import (
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type ActionResult struct {
	Action       model.Action
	Delta        model.Resources
	LocationName string
}

func ApplyAction(s *model.State, action model.Action) ActionResult {
	// TODO:
	// delta impacted by weather
	// build velocity gain lessens per new teammate
	// calculate game winning or losing
	out := ActionResult{Action: action}
	teamSize := len(s.Party.Members)
	switch action {
	case model.ActionTravel:
		s.CurrentLocation += 1
		out.LocationName = s.Route[s.CurrentLocation].Name
		out.Delta = model.Resources{
			Cash:   -150 * teamSize,
			Coffee: -2 * teamSize,
			Morale: -5,
		}
	case model.ActionRest:
		out.Delta = model.Resources{
			Cash:   -300 * teamSize,
			Coffee: 10,
			Morale: 10,
		}
	case model.ActionBuild:
		out.Delta = model.Resources{
			Coffee:    -4 * teamSize,
			Morale:    -5,
			Readiness: 3 * teamSize,
		}
	case model.ActionMarket:
		out.Delta = model.Resources{
			Cash:   -2000,
			Coffee: -2,
			Hype:   10,
		}
	default:
		panic("attempted to apply undefined action")
	}

	s.Resources.Add(out.Delta)
	s.Day += 1
	return out
}
