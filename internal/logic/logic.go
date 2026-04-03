package logic

import (
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type RNG interface {
	IntN(int) int
}

type Result struct {
	Delta model.Resources
	// CurrentLocation will always >= 1 after a Travel action, so we're reusing
	// default value 0 as sentinel value for "did not travel".
	CurrentLocation int
}

type ActionResult struct {
	Result
	Weather      model.WeatherKind
	WeatherDelta model.Resources
}

type EventResult struct {
	Result
}

type Context struct {
	Action model.Action
}

type Effect func(*model.State, Context) Change

type Change struct {
	Delta        model.Resources
	AdvanceDay   bool
	MoveLocation int
}

func (c Change) Apply(s *model.State) Result {
	out := Result{
		Delta: c.Delta,
	}
	if c.MoveLocation != 0 {
		s.CurrentLocation += c.MoveLocation
		out.CurrentLocation = s.CurrentLocation
	}
	s.Resources.AddClamped(c.Delta)
	if c.AdvanceDay {
		s.Day++
	}
	return out
}

func mergeChanges(changes ...Change) Change {
	var out Change
	for _, change := range changes {
		out.Delta.Add(change.Delta)
		out.MoveLocation += change.MoveLocation
		out.AdvanceDay = out.AdvanceDay || change.AdvanceDay
	}
	return out
}

func applyChanges(s *model.State, changes ...Change) Result {
	change := mergeChanges(changes...)
	out := change.Apply(s)
	updateMetaChanges(s)
	return out
}

// Generic name because this is the natural place to keep track of all
// side effects
func updateMetaChanges(s *model.State) {
	if s.Resources.Coffee == 0 {
		s.NoCoffeeDayCount++
	} else {
		s.NoCoffeeDayCount = 0
	}
}

func ApplyActionEffects(
	s *model.State,
	action model.Action,
	actionEffect Effect,
	weatherEffect Effect,
	_ RNG,
) ActionResult {
	if actionEffect == nil {
		panic("missing action effect")
	}

	ctx := Context{
		Action: action,
	}
	actionChange := actionEffect(s, ctx)
	weatherChange := Change{}
	if weatherEffect != nil {
		weatherChange = weatherEffect(s, ctx)
	}
	res := applyChanges(s, actionChange, weatherChange)
	return ActionResult{
		Result:       res,
		Weather:      s.Weather,
		WeatherDelta: weatherChange.Delta,
	}
}

func ApplyEventChoiceEffect(s *model.State, effect Effect) EventResult {
	if effect == nil {
		panic("missing event choice effect")
	}

	ctx := Context{}
	change := effect(s, ctx)
	return EventResult{Result: change.Apply(s)}
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

func resolveFinalPitch(s *model.State, rng RNG) Ending {
	finalPitchRoll := rng.IntN(100)
	offerThreshold := (s.Resources.Product + s.Resources.Hype/2)
	if finalPitchRoll < offerThreshold {
		return EndingTogether
	}
	return EndingNoOffer
}

func ResolveFinalEnding(s *model.State, rng RNG) Ending {
	ending := EvaluateEnding(s)
	if ending != EndingNone {
		return ending
	}
	return resolveFinalPitch(s, rng)
}
