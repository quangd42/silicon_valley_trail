package content

import (
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func TestLoad(t *testing.T) {
	cont := Load()
	if cont.Intro == "" {
		t.Fatal("missing intro copy")
	}
	if cont.Ending == "" {
		t.Fatal("missing ending copy")
	}
	if len(cont.Route) < 10 {
		t.Fatal("route is shorter than 10 locations")
	}
	if len(cont.Actions) != int(model.ActionCount) {
		t.Fatalf("count of available actions: %d\n count of actions with copy: %d", model.ActionCount, len(cont.Actions))
	}
	for action := range model.ActionCount {
		cop, ok := cont.Actions[action]
		if !ok {
			t.Fatalf("missing copy for action %d", action)
		}
		if cop.Desc == "" {
			t.Fatalf("missing action desc for action %d", action)
		}
		if len(cop.Narrative) < 1 {
			t.Fatalf("missing action narrative for action %d", action)
		}
	}
}
