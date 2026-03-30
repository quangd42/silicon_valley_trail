package model

import "testing"

func TestResourcesAdd(t *testing.T) {
	tests := []struct {
		name    string
		initial Resources
		delta   Resources
		want    Resources
	}{
		{
			name: "adds positive deltas",
			initial: Resources{
				Cash:    100,
				Morale:  50,
				Coffee:  2,
				Hype:    10,
				Product: 20,
			},
			delta: Resources{
				Cash:    25,
				Morale:  5,
				Coffee:  3,
				Hype:    8,
				Product: 12,
			},
			want: Resources{
				Cash:    125,
				Morale:  55,
				Coffee:  5,
				Hype:    18,
				Product: 32,
			},
		},
		{
			name: "applies mixed positive and negative deltas",
			initial: Resources{
				Cash:    500,
				Morale:  40,
				Coffee:  6,
				Hype:    20,
				Product: 30,
			},
			delta: Resources{
				Cash:    -200,
				Morale:  10,
				Coffee:  -2,
				Hype:    0,
				Product: 15,
			},
			want: Resources{
				Cash:    300,
				Morale:  50,
				Coffee:  4,
				Hype:    20,
				Product: 45,
			},
		},
		{
			name: "clamps negative totals at zero",
			initial: Resources{
				Cash:    100,
				Morale:  5,
				Coffee:  1,
				Hype:    3,
				Product: 7,
			},
			delta: Resources{
				Cash:    -150,
				Morale:  -10,
				Coffee:  -4,
				Hype:    -9,
				Product: -20,
			},
			want: Resources{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.initial
			got.Add(tt.delta)

			if got != tt.want {
				t.Fatalf("Resources.Add() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
