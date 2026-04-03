package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/view"
)

func TestRenderPrompt(t *testing.T) {
	def := gamedef.Load()
	promptView := view.DayPrompt(def)

	tests := []struct {
		name             string
		input            string
		wantChoice       model.PromptChoice
		wantSubstrings   []string
		extraPromptCount int
	}{
		{
			name:  "select first action",
			input: "1\n",
			wantChoice: model.PromptChoice{
				Kind:   model.ChoiceAction,
				Action: model.ActionTravel,
			},
			wantSubstrings: []string{
				"1. Travel to the next location (costs cash, coffee, and morale)\n",
				"5. Save Game\n",
				"Enter choice (1-6): ",
			},
		},
		{
			name:  "select last action",
			input: "4\n",
			wantChoice: model.PromptChoice{
				Kind:   model.ChoiceAction,
				Action: model.ActionMarket,
			},
			wantSubstrings: []string{
				"4. Marketing push (increase hype, costs a lot of cash and some coffee)\n",
				"Enter choice (1-6): ",
			},
		},
		{
			name:  "select first control",
			input: "5\n",
			wantChoice: model.PromptChoice{
				Kind:    model.ChoiceControl,
				Control: model.ControlSave,
			},
			wantSubstrings: []string{
				"5. Save Game\n",
				"6. Quit to Menu\n",
				"Enter choice (1-6): ",
			},
		},
		{
			name:  "invalid input retries before selecting control",
			input: "abc\n9\n6\n",
			wantChoice: model.PromptChoice{
				Kind:    model.ChoiceControl,
				Control: model.ControlQuitToMenu,
			},
			wantSubstrings: []string{
				"Invalid input. ",
			},
			extraPromptCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var b bytes.Buffer
			term := NewTerminal(reader, &b)

			gotChoice := term.RenderPrompt(promptView)
			if gotChoice != tt.wantChoice {
				t.Fatalf("want choice %#v, got choice %#v", tt.wantChoice, gotChoice)
			}

			gotDisplay := b.String()
			for _, want := range tt.wantSubstrings {
				if !strings.Contains(gotDisplay, want) {
					t.Fatalf("display missing %q in %q", want, gotDisplay)
				}
			}
			promptCount := strings.Count(gotDisplay, "Enter choice (1-6): ")
			if promptCount != tt.extraPromptCount+1 {
				t.Fatalf("prompt count = %d, want %d in %q", promptCount, tt.extraPromptCount+1, gotDisplay)
			}
		})
	}
}

func TestWrapText(t *testing.T) {
	t.Run("wraps long lines without breaking words", func(t *testing.T) {
		got := wrapText("A startup advisor tells you to ship less, explain more, and stop naming every button after a feeling.", 40)
		want := "A startup advisor tells you to ship\nless, explain more, and stop naming\nevery button after a feeling."
		if got != want {
			t.Fatalf("wrapText() = %q, want %q", got, want)
		}
	})
}
