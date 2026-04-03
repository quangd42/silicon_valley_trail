package save

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func TestJSONSaver(t *testing.T) {
	t.Run("save round trip", func(t *testing.T) {
		savePath := filepath.Join(t.TempDir(), "save.json")
		saver := NewJSONSaver(savePath)

		state := model.NewState(nil, nil, nil)
		state.Day = 4
		state.CurrentLocation = 3
		state.Resources = model.Resources{
			Cash:    8_250,
			Morale:  65,
			Coffee:  14,
			Hype:    22,
			Product: 31,
		}
		state.Party = model.Party{
			Members: []model.PartyMember{
				{Name: "You"},
				{Name: "Pete"},
				{Name: "Mina"},
			},
		}
		state.Weather = model.WeatherFog
		state.NoCoffeeDayCount = 2
		state.EventPools = model.EventPools{
			Main: []string{"helpful-founder"},
			Weather: map[model.WeatherKind][]string{
				model.WeatherClear: {"sunnydemo"},
				model.WeatherRainy: {"rainpitch"},
			},
		}
		state.CurrentEvent = "rainpitch"

		if err := saver.Save(state); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		info, err := os.Stat(savePath)
		if err != nil {
			t.Fatalf("saved file stat error = %v", err)
		}
		if info.Size() == 0 {
			t.Fatal("saved file is empty")
		}

		var loaded model.State
		if err := saver.Load(&loaded); err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if !reflect.DeepEqual(loaded, *state) {
			t.Fatalf("loaded state = %#v, want %#v", loaded, *state)
		}
	})

	t.Run("load corrupted save", func(t *testing.T) {
		savePath := filepath.Join(t.TempDir(), "save.json")
		if err := os.WriteFile(savePath, []byte("{not-json"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		saver := NewJSONSaver(savePath)
		initial := model.State{Day: 7}

		err := saver.Load(&initial)
		if !errors.Is(err, ErrSaveCorrupted) {
			t.Fatalf("Load() error = %v, want %v", err, ErrSaveCorrupted)
		}
		if initial.Day != 7 {
			t.Fatalf("Load() mutated state on error: got day %d, want 7", initial.Day)
		}
	})

	t.Run("correctly overwrites previous save", func(t *testing.T) {
		savePath := filepath.Join(t.TempDir(), "save.json")
		saver := NewJSONSaver(savePath)

		first := model.NewState(nil, nil, nil)
		first.Day = 1
		first.CurrentLocation = 1
		first.Resources = model.Resources{
			Cash:    9_500,
			Morale:  90,
			Coffee:  25,
			Hype:    15,
			Product: 23,
		}
		first.EventPools = model.EventPools{
			Main: []string{"first-event"},
			Weather: map[model.WeatherKind][]string{
				model.WeatherClear: {"clear-first"},
			},
		}
		first.CurrentEvent = "first-event"

		second := model.NewState(nil, nil, nil)
		second.Day = 8
		second.CurrentLocation = 4
		second.Resources = model.Resources{
			Cash:    5_200,
			Morale:  47,
			Coffee:  11,
			Hype:    34,
			Product: 41,
		}
		second.Weather = model.WeatherRainy
		second.NoCoffeeDayCount = 1
		second.EventPools = model.EventPools{
			Main: []string{"second-event"},
			Weather: map[model.WeatherKind][]string{
				model.WeatherRainy: {"rain-second"},
			},
		}
		second.CurrentEvent = "rain-second"

		if err := saver.Save(first); err != nil {
			t.Fatalf("first Save() error = %v", err)
		}
		if err := saver.Save(second); err != nil {
			t.Fatalf("second Save() error = %v", err)
		}

		var loaded model.State
		if err := saver.Load(&loaded); err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if !reflect.DeepEqual(loaded, *second) {
			t.Fatalf("loaded state = %#v, want %#v", loaded, *second)
		}
	})
}
