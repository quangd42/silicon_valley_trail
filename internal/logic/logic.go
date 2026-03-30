package logic

import (
	"math/rand"

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
			Coffee:  -4 * teamSize,
			Morale:  -5,
			Product: 3 * teamSize,
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
	if s.Resources.Coffee == 0 {
		s.NoCoffeeDayCount++
	} else if s.NoCoffeeDayCount > 0 {
		s.NoCoffeeDayCount = 0
	}
	s.Day += 1
	return out
}

type Ending int

const (
	EndingNone Ending = iota
	EndingNoCoffee
	EndingNoCash
	EndingNoOffer
	EndingTogether
	EndingAlone
)

func EvaluateEnding(s *model.State) Ending {
	if s.Resources.Cash == 0 {
		return EndingNoCash
	}
	if s.NoCoffeeDayCount > 1 {
		return EndingNoCoffee
	}
	return EndingNone
}

func resolveFinalPitch(s *model.State, finalPitchRoll int) Ending {
	offerThreshold := (s.Resources.Product + s.Resources.Hype/2)
	if finalPitchRoll < offerThreshold {
		return EndingAlone
	}
	return EndingNoOffer
}

func ResolveFinalEnding(s *model.State) Ending {
	ending := EvaluateEnding(s)
	if ending != EndingNone {
		return ending
	}
	return resolveFinalPitch(s, rand.Intn(100))
}
