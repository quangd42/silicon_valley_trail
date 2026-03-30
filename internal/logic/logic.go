package logic

import (
	"fmt"
	"log"

	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type ActionResult struct {
	Msg string
}

func ApplyAction(s *model.State, cont *content.Content, a model.Action) ActionResult {
	// TODO: Render game story messages based on Location
	// delta calculated based on Location
	// delta impacted by weather
	// calculate game winning or losing
	var delta actionDelta
	var msgArgs []any
	res := ActionResult{}
	teamSize := len(s.Party.Members)
	switch a {
	case model.ActionTravel:
		s.CurrentLocation += 1
		delta = actionDelta{
			Cash:   -100 * teamSize,
			Coffee: -2 * teamSize,
			Morale: -5,
		}
		msgArgs = []any{s.Route[s.CurrentLocation].Name}
	case model.ActionRest:
		delta = actionDelta{
			Cash:   -200 * teamSize,
			Coffee: 10,
			Morale: 10,
		}
	case model.ActionBuild:
		delta = actionDelta{
			Coffee:    -3 * teamSize,
			Morale:    -5,
			Readiness: 10,
		}
		msgArgs = []any{delta.Readiness, -delta.Coffee, -delta.Morale}
	case model.ActionMarket:
		delta = actionDelta{
			Cash: -2000,
			Hype: 10,
		}
		msgArgs = []any{-delta.Cash}
	}

	actionCopy, ok := cont.Actions[a]
	if !ok {
		log.Fatalf("missing copy for action id %d", a)
	}
	res.Msg = fmt.Sprintf(actionCopy.Result, msgArgs...)
	applyActionDelta(s, delta)
	s.Day += 1
	return res
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
