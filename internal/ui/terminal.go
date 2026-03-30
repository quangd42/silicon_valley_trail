package ui

// Implementation note: ALWAYS flush at the end of public `Render*` method!

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/view"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	thickSep = []byte(strings.Repeat("=", 80))
	thinSep  = []byte(strings.Repeat("-", 80))
)

type Terminal struct {
	in  *bufio.Reader
	out *bufio.Writer
	fmt *message.Printer
}

func NewTerminal(in io.Reader, out io.Writer) *Terminal {
	return &Terminal{
		in:  bufio.NewReader(in),
		out: bufio.NewWriter(out),
		fmt: message.NewPrinter(language.English),
	}
}

func (t *Terminal) RenderIntro(v view.IntroView) {
	t.RenderInfo(string(v))
}

func (t *Terminal) RenderDay(v view.DayView) {
	t.renderThickSep()
	t.fmt.Fprintf(t.out, "Day %d | %s\n", v.Day, v.Location.Name)
	t.fmt.Fprintf(t.out, "%s\n", v.Location.Desc)
	t.renderThickSep()
	r := v.Resources
	t.fmt.Fprintf(t.out, "Cash: $%d | Morale: %d%% | Coffee: %d\n", r.Cash, r.Morale, r.Coffee)
	t.fmt.Fprintf(t.out, "Hype: %d%% | Readiness: %d%%\n", r.Hype, r.Readiness)
	t.fmt.Fprintf(t.out, "Progress: %d%% to San Francisco\n", v.Progress)
	t.renderThickSep()
	t.fmt.Fprintf(t.out, "Weather: %s\n\n", v.Weather)
	t.renderThinSep()
	t.out.WriteString("What will you do?\n")
	t.renderThinSep()
	t.out.Flush()
}

// renderActions is typically
func (t *Terminal) renderActions(v []view.ActionView) {
	t.out.WriteString("Actions:\n")
	for i, action := range v {
		t.fmt.Fprintf(t.out, "%d. %s\n", i+1, action.Desc)
	}
	t.out.Write([]byte{'\n'})
}

func (t *Terminal) renderControls(v []model.Control, start int) {
	t.out.WriteString("Controls:\n")
	for i, control := range v {
		t.fmt.Fprintf(t.out, "%d. %s\n", i+start+1, control)
	}
	t.out.Write([]byte{'\n'})
}

// PromptChoice is the result type of `PromptSelection()`. It is a poor man's tagged union
// to distinguish if the user has chosen an in-game choice or a game session control action,
// as we unfortunately have to mix those two in this CLI representation.
// When `Kind` = true, `Action` is set, otherwise `Control` is set. Accessing the unset
// field does not panic, simply returns the default value.
type PromptChoice struct {
	Kind    bool // true = Action, false = Control
	Action  model.Action
	Control model.Control
}

func (t *Terminal) PromptSelection(v view.PromptView) PromptChoice {
	actionCount := len(v.Actions)
	controlCount := len(v.Controls)
	t.renderActions(v.Actions)
	t.renderControls(v.Controls, actionCount)
	for {
		t.fmt.Fprintf(t.out, "Enter choice (1-%d): ", actionCount+controlCount)
		t.out.Flush()
		input, err := t.in.ReadString('\n')
		if err != nil {
			// This often means Ctrl-C or some serious unrecoverable error, but we're
			// attempting to handle it anyway
			t.out.WriteByte('\n')
			continue
		}
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			t.out.WriteString("Invalid input. ")
			continue
		}
		switch {
		case choice >= 1 && choice <= actionCount:
			return PromptChoice{
				Kind:   true,
				Action: v.Actions[choice-1].Kind,
			}
		case choice > actionCount && choice <= actionCount+controlCount:
			return PromptChoice{
				Kind:    false,
				Control: v.Controls[choice-1-actionCount],
			}
		default:
			t.out.WriteString("Invalid input. ")
			continue

		}
	}
}

func (t *Terminal) renderImpact(location string, delta model.Resources) {
	if location != "" {
		t.fmt.Fprintf(t.out, "Arrived at %s!\n\n", location)
	}
	impacts := []string{}
	if delta.Cash != 0 {
		var sign byte
		if delta.Cash > 0 {
			sign = '+'
		} else {
			sign = '-'
			delta.Cash *= -1
		}
		impacts = append(impacts, t.fmt.Sprintf("Cash %c$%d", sign, delta.Cash))
	}
	if delta.Coffee != 0 {
		impacts = append(impacts, t.fmt.Sprintf("Coffee %+d", delta.Coffee))
	}
	if delta.Morale != 0 {
		impacts = append(impacts, t.fmt.Sprintf("Morale %+d%%", delta.Morale))
	}
	if delta.Hype != 0 {
		impacts = append(impacts, t.fmt.Sprintf("Hype %+d%%", delta.Hype))
	}
	if delta.Readiness != 0 {
		impacts = append(impacts, t.fmt.Sprintf("Product Readiness %+d%%", delta.Readiness))
	}
	if len(impacts) != 0 {
		t.fmt.Fprintf(t.out, "(%s)\n", strings.Join(impacts, ". "))
	}
}

func (t *Terminal) renderNarrative(v []string) {
	l := len(v)
	if l == 0 {
		return
	}
	for _, line := range v[:l-1] {
		t.out.WriteString(line)
		t.waitForEnter()
		t.clearLine()
	}
	t.out.WriteString(v[l-1])
	t.out.Flush()
}

func (t *Terminal) RenderActionResult(v view.ActionResultView) {
	t.renderThinSep()
	t.renderNarrative(v.Narative)
	t.renderImpact(v.LocationName, v.Delta)
	t.renderThinSep()
	t.waitForEnter()
	t.clearScreen()
}

func (t *Terminal) RenderInfo(msg string) {
	t.renderThinSep()
	t.out.WriteString(msg)
	t.renderThinSep()
	t.waitForEnter()
	t.clearScreen()
}

func (t *Terminal) RenderEnding(v view.EndingView) {
	t.fmt.Fprint(t.out, v)
	t.out.Flush()
}

func (t *Terminal) waitForEnter() {
	t.out.WriteString("Press Enter to continue...")
	t.out.Flush()
	t.in.ReadString('\n')
}

// print escape sequence to clear terminal screen (\033[2J) and move cursor to top left (\033[H)
func (t *Terminal) clearScreen() {
	t.out.WriteString("\033[2J\033[H")
}

// print escape sequence to move cursor up (\033[A) and clear entire line (\033[2K)
func (t *Terminal) clearLine() {
	t.out.WriteString("\033[A\033[2K")
}

func (t *Terminal) renderThickSep() {
	t.out.Write(thickSep)
	t.out.Write([]byte{'\n'})
}

func (t *Terminal) renderThinSep() {
	t.out.Write(thinSep)
	t.out.Write([]byte{'\n'})
}
