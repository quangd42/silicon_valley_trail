package logic

import (
	"reflect"
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func TestApplyAction(t *testing.T) {
	tests := []struct {
		name      string
		state     model.State
		action    model.Action
		want      ActionResult
		wantState model.State
	}{
		{
			name: "travel advances location and spends resources based on party size",
			state: model.State{
				Day: 3,
				Route: []model.Location{
					{Name: "San Jose"},
					{Name: "Santa Clara"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      10_000,
					Morale:    100,
					Coffee:    30,
					Hype:      10,
					Readiness: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
			action: model.ActionTravel,
			want: ActionResult{
				Action:       model.ActionTravel,
				Delta:        model.Resources{Cash: -300, Coffee: -4, Morale: -5},
				LocationName: "Santa Clara",
			},
			wantState: model.State{
				Day: 4,
				Route: []model.Location{
					{Name: "San Jose"},
					{Name: "Santa Clara"},
				},
				CurrentLocation: 1,
				Resources: model.Resources{
					Cash:      9_700,
					Morale:    95,
					Coffee:    26,
					Hype:      10,
					Readiness: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
		},
		{
			name: "rest restores morale and coffee while costing cash",
			state: model.State{
				Day: 1,
				Route: []model.Location{
					{Name: "San Jose"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      1_000,
					Morale:    40,
					Coffee:    2,
					Hype:      10,
					Readiness: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}, {Name: "Mina"}},
				},
			},
			action: model.ActionRest,
			want: ActionResult{
				Action: model.ActionRest,
				Delta:  model.Resources{Cash: -900, Coffee: 10, Morale: 10},
			},
			wantState: model.State{
				Day: 2,
				Route: []model.Location{
					{Name: "San Jose"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      100,
					Morale:    50,
					Coffee:    12,
					Hype:      10,
					Readiness: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}, {Name: "Mina"}},
				},
			},
		},
		{
			name: "build improves readiness and consumes morale and coffee",
			state: model.State{
				Day: 5,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      4_000,
					Morale:    60,
					Coffee:    9,
					Hype:      12,
					Readiness: 30,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
			action: model.ActionBuild,
			want: ActionResult{
				Action: model.ActionBuild,
				Delta:  model.Resources{Coffee: -8, Morale: -5, Readiness: 6},
			},
			wantState: model.State{
				Day: 6,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      4_000,
					Morale:    55,
					Coffee:    1,
					Hype:      12,
					Readiness: 36,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
		},
		{
			name: "market increases hype and clamps cash at zero",
			state: model.State{
				Day: 0,
				Route: []model.Location{
					{Name: "Palo Alto"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      1_500,
					Morale:    80,
					Coffee:    7,
					Hype:      15,
					Readiness: 25,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
			action: model.ActionMarket,
			want: ActionResult{
				Action: model.ActionMarket,
				Delta:  model.Resources{Cash: -2000, Coffee: -2, Hype: 10},
			},
			wantState: model.State{
				Day: 1,
				Route: []model.Location{
					{Name: "Palo Alto"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:      0,
					Morale:    80,
					Coffee:    5,
					Hype:      25,
					Readiness: 25,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ApplyAction(&tc.state, tc.action)

			if got != tc.want {
				t.Fatalf("ApplyAction() result = %#v, want %#v", got, tc.want)
			}
			if !reflect.DeepEqual(tc.state, tc.wantState) {
				t.Fatalf("ApplyAction() state = %#v, want %#v", tc.state, tc.wantState)
			}
		})
	}
}
