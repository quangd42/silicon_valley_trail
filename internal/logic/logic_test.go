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
					Cash:    10_000,
					Morale:  100,
					Coffee:  30,
					Hype:    10,
					Product: 20,
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
					Cash:    9_700,
					Morale:  95,
					Coffee:  26,
					Hype:    10,
					Product: 20,
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
					Cash:    1_000,
					Morale:  40,
					Coffee:  2,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}, {Name: "Mina"}},
				},
			},
			action: model.ActionRest,
			want: ActionResult{
				Action: model.ActionRest,
				Delta:  model.Resources{Cash: -900, Coffee: 12, Morale: 20},
			},
			wantState: model.State{
				Day: 2,
				Route: []model.Location{
					{Name: "San Jose"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    100,
					Morale:  60,
					Coffee:  14,
					Hype:    10,
					Product: 20,
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
					Cash:    4_000,
					Morale:  60,
					Coffee:  9,
					Hype:    12,
					Product: 30,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
			action: model.ActionBuild,
			want: ActionResult{
				Action: model.ActionBuild,
				Delta:  model.Resources{Coffee: -8, Morale: -15, Product: 3},
			},
			wantState: model.State{
				Day: 6,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  45,
					Coffee:  1,
					Hype:    12,
					Product: 33,
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
					Cash:    1_500,
					Morale:  80,
					Coffee:  7,
					Hype:    15,
					Product: 25,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
			action: model.ActionMarket,
			want: ActionResult{
				Action: model.ActionMarket,
				Delta:  model.Resources{Cash: -2000, Coffee: -4, Hype: 10},
			},
			wantState: model.State{
				Day: 1,
				Route: []model.Location{
					{Name: "Palo Alto"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    0,
					Morale:  80,
					Coffee:  3,
					Hype:    25,
					Product: 25,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
			},
		},
		{
			name: "cloudy weather gives rest morale+",
			state: model.State{
				Day: 2,
				Route: []model.Location{
					{Name: "San Jose"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    2_000,
					Morale:  50,
					Coffee:  5,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherCloudy,
			},
			action: model.ActionRest,
			want: ActionResult{
				Action:       model.ActionRest,
				Delta:        model.Resources{Cash: -600, Coffee: 12, Morale: 24},
				WeatherDelta: model.Resources{Morale: 4},
				Weather:      model.WeatherCloudy,
			},
			wantState: model.State{
				Day: 3,
				Route: []model.Location{
					{Name: "San Jose"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    1_400,
					Morale:  74,
					Coffee:  17,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherCloudy,
			},
		},
		{
			name: "fog weather gives build rate-",
			state: model.State{
				Day: 6,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  100,
					Coffee:  12,
					Hype:    12,
					Product: 30,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherFog,
			},
			action: model.ActionBuild,
			want: ActionResult{
				Action:       model.ActionBuild,
				Delta:        model.Resources{Coffee: -8, Morale: -15, Product: 5},
				WeatherDelta: model.Resources{Product: -1},
				Weather:      model.WeatherFog,
			},
			wantState: model.State{
				Day: 7,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  85,
					Coffee:  4,
					Hype:    12,
					Product: 35,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherFog,
			},
		},
		{
			name: "rainy weather gives travel morale- and coffee use+",
			state: model.State{
				Day: 0,
				Route: []model.Location{
					{Name: "San Jose"},
					{Name: "Santa Clara"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    10_000,
					Morale:  100,
					Coffee:  30,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherRainy,
			},
			action: model.ActionTravel,
			want: ActionResult{
				Action:       model.ActionTravel,
				Delta:        model.Resources{Cash: -300, Coffee: -6, Morale: -8},
				WeatherDelta: model.Resources{Coffee: -2, Morale: -3},
				Weather:      model.WeatherRainy,
				LocationName: "Santa Clara",
			},
			wantState: model.State{
				Day: 1,
				Route: []model.Location{
					{Name: "San Jose"},
					{Name: "Santa Clara"},
				},
				CurrentLocation: 1,
				Resources: model.Resources{
					Cash:    9_700,
					Morale:  92,
					Coffee:  24,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherRainy,
			},
		},
		{
			name: "clear weather gives build rate+ and morale+",
			state: model.State{
				Day: 4,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  100,
					Coffee:  12,
					Hype:    12,
					Product: 30,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherClear,
			},
			action: model.ActionBuild,
			want: ActionResult{
				Action:       model.ActionBuild,
				Delta:        model.Resources{Coffee: -8, Morale: -12, Product: 8},
				WeatherDelta: model.Resources{Morale: 3, Product: 2},
				Weather:      model.WeatherClear,
			},
			wantState: model.State{
				Day: 5,
				Route: []model.Location{
					{Name: "Mountain View"},
				},
				CurrentLocation: 0,
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  88,
					Coffee:  4,
					Hype:    12,
					Product: 38,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				Weather: model.WeatherClear,
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

func TestApplyAction_NoCoffeeDayCount(t *testing.T) {
	tests := []struct {
		name                 string
		state                model.State
		action               model.Action
		wantNoCoffeeDayCount int
	}{
		{
			name: "increments no-coffee counter when coffee reaches zero",
			state: model.State{
				Route: []model.Location{{Name: "San Jose"}},
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  60,
					Coffee:  4,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}},
				},
			},
			action:               model.ActionBuild,
			wantNoCoffeeDayCount: 1,
		},
		{
			name: "resets no-coffee counter when coffee is restored",
			state: model.State{
				Route: []model.Location{{Name: "San Jose"}},
				Resources: model.Resources{
					Cash:    4_000,
					Morale:  60,
					Coffee:  0,
					Hype:    10,
					Product: 20,
				},
				Party: model.Party{
					Members: []model.PartyMember{{Name: "You"}, {Name: "Pete"}},
				},
				NoCoffeeDayCount: 1,
			},
			action:               model.ActionRest,
			wantNoCoffeeDayCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ApplyAction(&tc.state, tc.action)
			if tc.state.NoCoffeeDayCount != tc.wantNoCoffeeDayCount {
				t.Fatalf("NoCoffeeDayCount = %d, want %d", tc.state.NoCoffeeDayCount, tc.wantNoCoffeeDayCount)
			}
		})
	}
}

func TestEvaluateEnding(t *testing.T) {
	tests := []struct {
		name  string
		state model.State
		want  Ending
	}{
		{
			name: "returns none while run is still valid",
			state: model.State{
				Resources:        model.Resources{Cash: 100, Coffee: 5},
				NoCoffeeDayCount: 1,
			},
			want: EndingNone,
		},
		{
			name: "returns no coffee after two zero-coffee days",
			state: model.State{
				Resources:        model.Resources{Cash: 100, Coffee: 0},
				NoCoffeeDayCount: 2,
			},
			want: EndingNoCoffee,
		},
		{
			name: "returns no cash when cash is depleted",
			state: model.State{
				Resources:        model.Resources{Cash: 0, Coffee: 5},
				NoCoffeeDayCount: 0,
			},
			want: EndingNoCash,
		},
		{
			name: "no cash takes precedence over no coffee",
			state: model.State{
				Resources:        model.Resources{Cash: 0, Coffee: 0},
				NoCoffeeDayCount: 2,
			},
			want: EndingNoCash,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := EvaluateEnding(&tc.state)
			if got != tc.want {
				t.Fatalf("EvaluateEnding() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResolveFinalPitch(t *testing.T) {
	tests := []struct {
		name  string
		state model.State
		roll  int
		want  Ending
	}{
		{
			name: "loses when threshold is zero",
			state: model.State{
				Resources: model.Resources{Product: 0, Hype: 0},
			},
			roll: 0,
			want: EndingNoOffer,
		},
		{
			name: "wins when roll is below threshold",
			state: model.State{
				Resources: model.Resources{Product: 40, Hype: 20},
			},
			roll: 49,
			want: EndingAlone,
		},
		{
			name: "loses when roll equals threshold",
			state: model.State{
				Resources: model.Resources{Product: 40, Hype: 20},
			},
			roll: 50,
			want: EndingNoOffer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveFinalPitch(&tc.state, tc.roll)
			if got != tc.want {
				t.Fatalf("resolveFinalPitch() = %v, want %v", got, tc.want)
			}
		})
	}
}
