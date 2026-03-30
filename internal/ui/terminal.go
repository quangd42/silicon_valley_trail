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
	t.out.Write([]byte("What will you do?\n"))
	t.renderThinSep()
	t.out.Flush()
}

// renderActions is typically
func (t *Terminal) renderActions(v []view.ActionView) {
	for i, action := range v {
		t.fmt.Fprintf(t.out, "%d. %s\n", i+1, action.Desc)
	}
	t.out.Write([]byte{'\n'})
}

func (t *Terminal) Prompt(v view.PromptView) model.Action {
	t.renderActions(v.Actions)
	for {
		t.fmt.Fprintf(t.out, "Enter choice (1-%d): ", len(v.Actions))
		t.out.Flush()
		input, err := t.in.ReadString('\n')
		if err != nil {
			// This often means Ctrl-C or some serious unrecoverable error, but we're
			// attempting to handle it anyway
			t.out.Write([]byte{'\n'})
			continue
		}
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(v.Actions) {
			t.out.Write([]byte("Invalid input. "))
			continue
		}
		return v.Actions[choice-1].Kind
	}
}

func (t *Terminal) RenderInfo(msg string) {
	t.out.Write([]byte{'\n'})
	t.out.Write([]byte(msg))
	t.out.Write([]byte("Press Enter to continue..."))
	t.out.Flush()
	t.in.ReadString('\n')
	t.clearScreen()
}

func (t *Terminal) RenderEnding(v view.EndingView) {
	t.fmt.Fprint(t.out, v)
	t.out.Flush()
}

// print escape sequence to clear terminal screen (\033[2J) and move cursor to top left (\033[H)
func (t *Terminal) clearScreen() {
	t.out.Write([]byte("\033[2J\033[H"))
}

func (t *Terminal) renderThickSep() {
	t.out.Write(thickSep)
	t.out.Write([]byte{'\n'})
}

func (t *Terminal) renderThinSep() {
	t.out.Write(thinSep)
	t.out.Write([]byte{'\n'})
}
