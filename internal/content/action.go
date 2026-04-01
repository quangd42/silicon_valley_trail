package content

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type ActionData struct {
	Desc      string
	Narrative Narrative
	Effect    logic.Effect
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
						Cash:   -150 * teamSize,
						Coffee: -2 * teamSize,
						Morale: -5,
					},
					MoveLocation: 1,
					AdvanceDay:   true,
				}
			},
		},
		model.ActionRest: {
			Desc:      "Rest and recover (restore morale and coffee, costs cash)",
			Narrative: []string{"You decided to take a break...", "...zZz...", "You feel refreshed. You're filled with determination."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Cash:   -300 * teamSize,
						Coffee: 12,
						Morale: 20,
					},
					AdvanceDay: true,
				}
			},
		},
		model.ActionBuild: {
			Desc:      "Work on product (increase product readiness, costs coffee and morale)",
			Narrative: []string{"You take on the next item on the roadmap...", "You're happy with the result, but everyone is tired..."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Coffee:  -4 * teamSize,
						Product: 3 * teamSize * s.Resources.Morale / 100,
						Morale:  -15,
					},
					AdvanceDay: true,
				}
			},
		},
		model.ActionMarket: {
			Desc:      "Marketing push (increase hype, costs cash and coffee)",
			Narrative: []string{"You launch a marketing campaign...", "Every \"debate\" on X is about your product."},
			Effect: func(s *model.State, _ logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				return logic.Change{
					Delta: model.Resources{
						Cash:   -2000,
						Coffee: -2 * teamSize,
						Hype:   10,
					},
					AdvanceDay: true,
				}
			},
		},
	}
}
