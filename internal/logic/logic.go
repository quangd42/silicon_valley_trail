package logic

import (
	"math/rand"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type ActionResult struct {
	Action       model.Action
	Delta        model.Resources
	Weather      model.WeatherKind
	WeatherDelta model.Resources
	// CurrentLocation will always >= 1 after a Travel action, so we're reusing
	// default value 0 as sentinel value for "did not travel".
	CurrentLocation int
}

func ApplyAction(s *model.State, action model.Action) ActionResult {
	// TODO: build velocity gain lessens per new teammate
	out := ActionResult{
		Action:  action,
		Weather: s.Weather,
	}
	teamSize := len(s.Party.Members)
	switch action {
	case model.ActionTravel:
		s.CurrentLocation += 1
		out.CurrentLocation = s.CurrentLocation
		out.Delta = model.Resources{
			Cash:   -150 * teamSize,
			Coffee: -2 * teamSize,
			Morale: -5,
		}
	case model.ActionRest:
		out.Delta = model.Resources{
			Cash:   -300 * teamSize,
			Coffee: 12,
			Morale: 20,
		}
	case model.ActionBuild:
		out.Delta = model.Resources{
			Coffee:  -4 * teamSize,
			Product: 3 * teamSize * s.Resources.Morale / 100,
			Morale:  -15,
		}
	case model.ActionMarket:
		out.Delta = model.Resources{
			Cash:   -2000,
			Coffee: -2 * teamSize,
			Hype:   10,
		}
	default:
		panic("attempted to apply undefined action")
	}

	out.WeatherDelta = computeWeatherImpact(s, action)
	out.Delta.Add(out.WeatherDelta)
	s.Resources.AddClamped(out.Delta)
	if s.Resources.Coffee == 0 {
		s.NoCoffeeDayCount++
	} else if s.NoCoffeeDayCount > 0 {
		s.NoCoffeeDayCount = 0
	}
	s.Day += 1
	return out
}

func computeWeatherImpact(s *model.State, action model.Action) model.Resources {
	teamSize := len(s.Party.Members)
	switch s.Weather {
	case model.WeatherClear:
		switch action {
		case model.ActionTravel:
			return model.Resources{Morale: 2}
		case model.ActionBuild:
			return model.Resources{
				Morale:  3,
				Product: teamSize * s.Resources.Morale / 100,
			}
		}
	case model.WeatherCloudy:
		switch action {
		case model.ActionRest:
			return model.Resources{Morale: 4}
		case model.ActionMarket:
			return model.Resources{Hype: 2}
		}
	case model.WeatherRainy:
		switch action {
		case model.ActionTravel:
			return model.Resources{
				Coffee: -teamSize,
				Morale: -3,
			}
		case model.ActionMarket:
			return model.Resources{Hype: -2}
		}
	case model.WeatherFog:
		switch action {
		case model.ActionTravel:
			return model.Resources{Morale: -2}
		case model.ActionBuild:
			return model.Resources{Product: -1}
		default:
			return model.Resources{}
		}
	}
	return model.Resources{}
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
