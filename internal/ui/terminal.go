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

func (t *Terminal) RenderMainMenu(v view.PromptView) PromptChoice {
	t.thickSep()
	t.out.WriteString("SILICON VALLEY TRAIL - Main Menu\n")
	t.thickSep()
	return t.PromptSelection(v)
}

func (t *Terminal) RenderIntro(v view.IntroView) {
	t.renderNarrative(v)
	t.thinSep()
	t.out.WriteString("Press Enter to begin your journey!\n")
	t.thinSep()
	t.out.Flush()
	t.in.ReadString('\n')
}

func (t *Terminal) RenderDay(v view.DayView) {
	t.thickSep()
	t.fmt.Fprintf(t.out, "Day %d | %s\n", v.Day, v.Location.Name)
	t.fmt.Fprintf(t.out, "%s\n", v.Location.Desc)
	t.thickSep()
	r := v.Resources
	t.fmt.Fprintf(t.out, "Cash: $%d | Morale: %d%% | Coffee: %d\n", r.Cash, r.Morale, r.Coffee)
	t.fmt.Fprintf(t.out, "Hype: %d%% | Readiness: %d%%\n", r.Hype, r.Readiness)
	t.fmt.Fprintf(t.out, "Progress: %d%% to San Francisco\n", v.Progress)
	t.thickSep()
	t.fmt.Fprintf(t.out, "Weather: %s\n\n", v.Weather)
	t.thinSep()
	t.out.WriteString("What will you do?\n")
	t.thinSep()
	t.out.Flush()
}

func (t *Terminal) renderActions(v []view.ActionView, label string) {
	t.out.WriteString(label)
	for i, action := range v {
		t.fmt.Fprintf(t.out, "%d. %s\n", i+1, action.Desc)
	}
	t.linefeed()
}

func (t *Terminal) renderControls(v []model.Control, label string, start int) {
	t.out.WriteString(label)
	for i, control := range v {
		t.fmt.Fprintf(t.out, "%d. %s\n", i+start+1, control)
	}
	t.linefeed()
}

// PromptChoice is the type of the result of `PromptSelection()`. It is a poor man's tagged
// union, to distinguish if the user has chosen an in-game action or a game session control,
// as we unfortunately have to mix those two choices in this CLI UI representation.
// When `Kind` = true, `Action` is set, otherwise `Control` is set. Accessing the unset
// field does not panic, simply returns the default (and wrong) value.
type PromptChoice struct {
	Kind    bool // true = Action, false = Control
	Action  model.Action
	Control model.Control
}

func (t *Terminal) PromptSelection(v view.PromptView) PromptChoice {
	actionCount := len(v.Actions)
	controlCount := len(v.Controls)
	if actionCount > 0 {
		t.renderActions(v.Actions, v.ActionsLabel)
	}
	if controlCount > 0 {
		t.renderControls(v.Controls, v.ControlsLabel, actionCount)
	}
	totalCount := actionCount + controlCount
	if totalCount == 0 {
		panic("prompt view has no choice to render")
	}
	for {
		t.fmt.Fprintf(t.out, "Enter choice (1-%d): ", totalCount)
		t.out.Flush()
		input, err := t.in.ReadString('\n')
		if err != nil {
			// This often means Ctrl-C or some serious unrecoverable error, but we're
			// attempting to handle it anyway
			t.linefeed()
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
	t.thinSep()
	t.renderNarrative(v.Narative)
	t.renderImpact(v.LocationName, v.Delta)
	t.thinSep()
	t.waitForEnter()
	t.ClearScreen()
}

func (t *Terminal) RenderInfo(msg string) {
	t.thinSep()
	t.out.WriteString(msg)
	t.linefeed()
	t.thinSep()
	t.waitForEnter()
	t.ClearScreen()
}

func (t *Terminal) RenderInfoNoWait(msg string) {
	t.thinSep()
	t.out.WriteString(msg)
	t.linefeed()
	t.thinSep()
	t.out.Flush()
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

// ClearScreen prints escape sequence to clear terminal screen (\033[2J) and move cursor to top left (\033[H)
func (t *Terminal) ClearScreen() {
	t.out.WriteString("\033[2J\033[H")
	t.out.Flush()
}

// print escape sequence to move cursor up (\033[A) and clear entire line (\033[2K)
func (t *Terminal) clearLine() {
	t.out.WriteString("\033[A\033[2K")
}

func (t *Terminal) thickSep() {
	t.out.Write(thickSep)
	t.linefeed()
}

func (t *Terminal) thinSep() {
	t.out.Write(thinSep)
	t.linefeed()
}

func (t *Terminal) linefeed() {
	t.out.WriteByte('\n')
}
