package ui

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
	out io.Writer
	fmt *message.Printer
}

func NewTerminal(in io.Reader, out io.Writer) *Terminal {
	return &Terminal{
		in:  bufio.NewReader(in),
		out: out,
		fmt: message.NewPrinter(language.English),
	}
}

func (t *Terminal) RenderIntro(v view.IntroView) {
	t.fmt.Fprint(t.out, v)
}

func (t *Terminal) RenderDay(v view.DayView) {
	// Really hacky way to clear the screen on each day:
	// print escape sequence to clear screen (\033[2J) and move cursor to top left (\033[H)
	t.out.Write([]byte("\033[2J\033[H"))
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
	t.out.Write([]byte("What will you do?\n"))
	t.renderThinSep()
}

func (t *Terminal) renderActions(actions []model.ActionKind) {
	for i, action := range actions {
		t.fmt.Fprintf(t.out, "%d. %s\n", i+1, action)
	}
	t.out.Write([]byte{'\n'})
}

func (t *Terminal) Prompt(actions []model.ActionKind) model.ActionKind {
	t.renderActions(actions)
	for {
		t.fmt.Fprintf(t.out, "Enter choice (1-%d): ", len(actions))
		input, err := t.in.ReadString('\n')
		if err != nil {
			// Only if for some reason the user forces input without typing '\n'
			t.out.Write([]byte{'\n'})
			continue
		}
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			t.out.Write([]byte("Invalid input. "))
			continue
		}
		if choice < 1 || choice > len(actions) {
			continue
		}
		return actions[choice-1]
	}
}

func (t *Terminal) RenderInfo(msg string) {
	t.out.Write([]byte{'\n'})
	t.out.Write([]byte(msg))
}

func (t *Terminal) RenderEnding(v view.EndingView) {
	t.fmt.Fprint(t.out, v)
}

func (t *Terminal) renderThickSep() {
	t.out.Write(thickSep)
	t.out.Write([]byte{'\n'})
}

func (t *Terminal) renderThinSep() {
	t.out.Write(thinSep)
	t.out.Write([]byte{'\n'})
}
