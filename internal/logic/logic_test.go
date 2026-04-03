package logic

import (
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func TestChangeApply(t *testing.T) {
	state := model.State{
		Day:             4,
		CurrentLocation: 1,
		Resources: model.Resources{
			Cash:    100,
			Morale:  10,
			Coffee:  2,
			Hype:    3,
			Product: 4,
		},
	}

	got := Change{
		Delta:        model.Resources{Cash: -150, Morale: 5, Coffee: -5, Hype: 2},
		MoveLocation: 1,
		AdvanceDay:   true,
	}.Apply(&state)

	want := Result{
		Delta:           model.Resources{Cash: -150, Morale: 5, Coffee: -5, Hype: 2},
		CurrentLocation: 2,
	}
	if got != want {
		t.Fatalf("Change.Apply() result = %#v, want = %#v", got, want)
	}
	if state.Day != 5 {
		t.Fatalf("Day = %d, want 5", state.Day)
	}
	if state.CurrentLocation != 2 {
		t.Fatalf("CurrentLocation = %d, want 2", state.CurrentLocation)
	}
	if state.Resources != (model.Resources{
		Cash:    0,
		Morale:  15,
		Coffee:  0,
		Hype:    5,
		Product: 4,
	}) {
		t.Fatalf("Resources = %#v", state.Resources)
	}
}

func TestMergeChanges(t *testing.T) {
	got := mergeChanges(
		Change{
			Delta:        model.Resources{Cash: -10, Morale: 1},
			MoveLocation: 1,
			AdvanceDay:   true,
		},
		Change{
			Delta:        model.Resources{Cash: 5, Coffee: -2},
			MoveLocation: 2,
		},
	)

	want := Change{
		Delta:        model.Resources{Cash: -5, Morale: 1, Coffee: -2},
		MoveLocation: 3,
		AdvanceDay:   true,
	}
	if got != want {
		t.Fatalf("mergeChanges() = %#v, want %#v", got, want)
	}
}

func TestApplyActionEffects(t *testing.T) {
	t.Run("merges resolved action and weather changes into one applied turn", func(t *testing.T) {
		state := model.State{
			Day:             1,
			CurrentLocation: 2,
			Resources:       model.Resources{Cash: 100, Morale: 10, Coffee: 5},
			Weather:         model.WeatherCloudy,
		}

		actionEffect := func(_ *model.State, _ Context) Change {
			return Change{
				Delta:        model.Resources{Cash: -30, Coffee: -2},
				MoveLocation: 1,
				AdvanceDay:   true,
			}
		}
		weatherEffect := func(s *model.State, _ Context) Change {
			return Change{
				Delta: model.Resources{Morale: s.Resources.Cash / 25},
			}
		}

		got := ApplyActionEffects(&state, model.ActionTravel, actionEffect, weatherEffect, nil)

		want := ActionResult{
			Result: Result{
				Delta:           model.Resources{Cash: -30, Coffee: -2, Morale: 4},
				CurrentLocation: 3,
			},
			Weather:      model.WeatherCloudy,
			WeatherDelta: model.Resources{Morale: 4},
		}
		if got != want {
			t.Fatalf("ApplyActionEffects() = %#v, want %#v", got, want)
		}
		if state.Day != 2 {
			t.Fatalf("Day = %d, want 2", state.Day)
		}
		if state.CurrentLocation != 3 {
			t.Fatalf("CurrentLocation = %d, want 3", state.CurrentLocation)
		}
		if state.Resources != (model.Resources{Cash: 70, Morale: 14, Coffee: 3}) {
			t.Fatalf("Resources = %#v", state.Resources)
		}
	})

	t.Run("increments streak when combined change leaves coffee at zero", func(t *testing.T) {
		state := model.State{
			Resources:        model.Resources{Coffee: 5},
			NoCoffeeDayCount: 1,
		}

		ApplyActionEffects(
			&state,
			model.ActionRest,
			func(_ *model.State, _ Context) Change {
				return Change{Delta: model.Resources{Coffee: -4}}
			},
			func(_ *model.State, _ Context) Change {
				return Change{Delta: model.Resources{Coffee: -1}}
			},
			nil,
		)

		if state.NoCoffeeDayCount != 2 {
			t.Fatalf("NoCoffeeDayCount = %d, want 2", state.NoCoffeeDayCount)
		}
	})

	t.Run("resets streak when coffee is restored", func(t *testing.T) {
		state := model.State{
			Resources:        model.Resources{Coffee: 0},
			NoCoffeeDayCount: 2,
		}

		ApplyActionEffects(
			&state,
			model.ActionRest,
			func(_ *model.State, _ Context) Change {
				return Change{Delta: model.Resources{Coffee: 3}}
			},
			nil,
			nil,
		)

		if state.NoCoffeeDayCount != 0 {
			t.Fatalf("NoCoffeeDayCount = %d, want 0", state.NoCoffeeDayCount)
		}
	})
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

type Rand struct {
	roll int
}

func (r Rand) IntN(_ int) int {
	return r.roll
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
			want: EndingTogether,
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
			got := resolveFinalPitch(&tc.state, Rand{roll: tc.roll})
			if got != tc.want {
				t.Fatalf("resolveFinalPitch() = %v, want %v", got, tc.want)
			}
		})
	}
}
