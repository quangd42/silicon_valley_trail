package logic

import (
	"fmt"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func ApplyAction(s *model.State, a model.ActionKind) string {
	// TODO: Render game story messages based on Location
	// delta calculated based on Location
	// delta impacted by weather
	// calculate game winning or losing
	var delta actionDelta
	var msg string
	switch a {
	case model.ActionTravel:
		s.CurrentLocation += 1
		delta = actionDelta{
			Cash:   -100 * len(s.Party.Members),
			Coffee: -(s.Resources.Coffee / 10),
			Morale: -(s.Resources.Morale / 20),
		}
		msg = fmt.Sprintf("Your team hit the road...\n\nArrived at %s!\n\n", s.Route[s.CurrentLocation].Name)
	case model.ActionRest:
		delta = actionDelta{
			Cash:   -200 * len(s.Party.Members),
			Coffee: s.Resources.Coffee / 20,
			Morale: s.Resources.Morale / 10,
		}
		msg = "You decided to take a break.\n\nYou're filled with determination.\n\n"
	case model.ActionBuild:
		delta = actionDelta{
			Coffee:    -(s.Resources.Coffee / 5),
			Morale:    -(s.Resources.Morale / 20),
			Readiness: 10,
		}
		msg = "You take on the next item on the roadmap...\n\nYou're happy with the result, but everyone is tired...\n\n"
	case model.ActionMarket:
		delta = actionDelta{
			Cash: -2000,
			Hype: s.Resources.Morale / 10,
		}
		msg = "You launch a marketing campaign...\n\nEvery startup in CA is talking about you...(Hype increased. Cost $2,000)\n\n"
	}
	applyActionDelta(s, delta)
	s.Day += 1
	return msg
}

type actionDelta = model.Resources

func applyActionDelta(s *model.State, delta actionDelta) {
	s.Resources.Cash += delta.Cash
	if s.Resources.Cash < 0 {
		s.Resources.Cash = 0
	}
	s.Resources.Coffee += delta.Coffee
	if s.Resources.Coffee < 0 {
		s.Resources.Coffee = 0
	}
	s.Resources.Morale += delta.Morale
	if s.Resources.Morale < 0 {
		s.Resources.Morale = 0
	}
	s.Resources.Hype += delta.Hype
	if s.Resources.Hype < 0 {
		s.Resources.Hype = 0
	}
	s.Resources.Readiness += delta.Readiness
	if s.Resources.Readiness < 0 {
		s.Resources.Readiness = 0
	}
}
