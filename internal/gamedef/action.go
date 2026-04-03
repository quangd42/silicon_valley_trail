package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type ActionData struct {
	Desc      string
	Narrative Narrative
	Effect    logic.Effect
}

func actionOrder() []model.Action {
	return []model.Action{
		model.ActionTravel,
		model.ActionRest,
		model.ActionBuild,
		model.ActionMarket,
	}
}

func actionData() map[model.Action]ActionData {
	return map[model.Action]ActionData{
		model.ActionTravel: {
			Desc:      "Travel to the next location (costs cash, coffee, and morale)",
			Narrative: []string{"Your team hit the road..."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Cash:   -175 * teamSize,
						Coffee: -2 * teamSize,
						Morale: -6,
					},
					MoveLocation: 1,
					AdvanceDay:   true,
				}
			},
		},
		model.ActionRest: {
			Desc:      "Rest and recover (restore morale and coffee, costs cash)",
			Narrative: []string{"You decided to take a break...\n\n...zZz...\n\nYou feel refreshed. You're filled with determination."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Cash:   -225 * teamSize,
						Coffee: 9,
						Morale: 18,
					},
					AdvanceDay: true,
				}
			},
		},
		model.ActionBuild: {
			Desc:      "Work on product (increase product readiness, costs coffee and morale)",
			Narrative: []string{"You take on the next item on the roadmap...\n\nYou're happy with the result, but everyone is tired..."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Coffee:  -3 * teamSize,
						Product: 4 * teamSize * s.Resources.Morale / 100,
						Morale:  -18,
					},
					AdvanceDay: true,
				}
			},
		},
		model.ActionMarket: {
			Desc:      "Marketing push (increase hype, costs a lot of cash and some coffee)",
			Narrative: []string{"You launch a marketing campaign...\n\nEvery \"debate\" on X is about your product."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Cash:   -850,
						Coffee: -1 * teamSize,
						Hype:   18,
					},
					AdvanceDay: true,
				}
			},
		},
	}
}
