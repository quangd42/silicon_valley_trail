package gamedef

import (
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func TestLoad(t *testing.T) {
	def := Load()
	if len(def.Intro) == 0 {
		t.Fatal("missing intro")
	}
	if len(def.Endings) == 0 {
		t.Fatal("missing ending definitions")
	}
	if len(def.Route) < 10 {
		t.Fatal("route is shorter than 10 locations")
	}
	if len(def.Actions) != int(model.ActionCount) {
		t.Fatalf("count of available actions: %d\ncount of defined actions: %d", model.ActionCount, len(def.Actions))
	}
	for action := range model.ActionCount {
		actionDef, ok := def.Actions[action]
		if !ok {
			t.Fatalf("missing definition for action %d", action)
		}
		if actionDef.Desc == "" {
			t.Fatalf("missing action desc for action %d", action)
		}
		if len(actionDef.Narrative) < 1 {
			t.Fatalf("missing action narrative for action %d", action)
		}
		if actionDef.Effect == nil {
			t.Fatalf("missing action effect for action %d", action)
		}
	}
	for weather := range model.WeatherKindCount {
		weatherDef, ok := def.Weather[weather]
		if !ok {
			t.Fatalf("missing definition for weather %d", weather)
		}
		if weatherDef.Desc == "" {
			t.Fatalf("missing weather desc for weather %d", weather)
		}
		if weather != model.WeatherUnknown && weatherDef.Effect == nil {
			t.Fatalf("missing weather effect for weather %d", weather)
		}
	}
}
