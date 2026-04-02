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

const terminalWidth = 80

var (
	thickSep = []byte(strings.Repeat("=", terminalWidth))
	thinSep  = []byte(strings.Repeat("-", terminalWidth))
	alertSep = []byte(strings.Repeat("!", terminalWidth))
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
	t.renderNarrative(v)
	t.thinSep()
	t.out.WriteString("Press Enter to begin your journey!\n")
	t.thinSep()
	t.out.Flush()
	t.in.ReadString('\n')
	t.ClearScreen()
}

func (t *Terminal) RenderDayInfo(v view.DayView) {
	t.thickSep()
	t.fmt.Fprintf(t.out, "Day %d | %s\n", v.Day, v.Location.Name)
	t.fmt.Fprintf(t.out, "%s\n", wrapText(v.Location.Desc, terminalWidth))
	t.thickSep()
	r := v.Resources
	t.fmt.Fprintf(t.out, "Cash: $%d | Morale: %d%% | Coffee: %d\n", r.Cash, r.Morale, r.Coffee)
	t.fmt.Fprintf(t.out, "Hype: %d%% | Product Readiness: %d%%\n", r.Hype, r.Product)
	t.fmt.Fprintf(t.out, "Progress: %d%% to San Francisco\n", v.Progress)
	t.thickSep()
	t.fmt.Fprintf(t.out, "Weather: %s\n%s\n", v.Weather, wrapText(v.WeatherImpact, terminalWidth))
	t.thinSep()
	t.out.WriteString("What will you do?\n")
	t.thinSep()
	t.out.Flush()
}

func (t *Terminal) RenderPrompt(v view.PromptView) model.PromptChoice {
	if v.Title != "" {
		t.thickSep()
		t.out.WriteString(v.Title)
		t.linefeed()
		t.thickSep()
	}
	choices := make([]model.PromptChoice, 0)
	for _, section := range v.Sections {
		if len(section.Items) == 0 {
			continue
		}
		if section.Label != "" {
			t.out.WriteString(section.Label)
			t.linefeed()
		}
		for _, item := range section.Items {
			choices = append(choices, item.Choice)
			t.fmt.Fprintf(t.out, "%d. %s\n", len(choices), item.Text)
		}
		t.linefeed()
	}
	choice := t.readPromptSelection(len(choices))
	return choices[choice-1]
}

func (t *Terminal) readPromptSelection(totalCount int) int {
	if totalCount == 0 {
		panic("no choice to render")
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
		if err != nil || choice < 1 || choice > totalCount {
			t.out.WriteString("Invalid input. ")
			continue
		}
		return choice
	}
}

func (t *Terminal) RenderActionResult(v view.ActionResultView) {
	t.thinSep()
	t.renderNarrative(v.Narrative)
	t.thinSep()
	t.renderImpact(v.LocationName, v.Delta)
	t.ClearScreen()
}

func (t *Terminal) renderImpact(location string, delta model.Resources) {
	if location != "" {
		t.fmt.Fprintf(t.out, "Arrived at %s!\n\n", location)
	}
	parts := []string{}
	if delta.Cash != 0 {
		var sign byte
		if delta.Cash > 0 {
			sign = '+'
		} else {
			sign = '-'
			delta.Cash *= -1
		}
		parts = append(parts, t.fmt.Sprintf("Cash %c$%d", sign, delta.Cash))
	}
	if delta.Coffee != 0 {
		parts = append(parts, t.fmt.Sprintf("Coffee %+d", delta.Coffee))
	}
	if delta.Morale != 0 {
		parts = append(parts, t.fmt.Sprintf("Morale %+d%%", delta.Morale))
	}
	if delta.Hype != 0 {
		parts = append(parts, t.fmt.Sprintf("Hype %+d%%", delta.Hype))
	}
	if delta.Product != 0 {
		parts = append(parts, t.fmt.Sprintf("Product %+d%%", delta.Product))
	}
	parts = append(parts, "Press Enter to continue...")
	t.waitForEnterMsg(t.fmt.Sprintf("(%s)", strings.Join(parts, ". ")))
}

func (t *Terminal) RenderEventInfo(v view.EventView) {
	t.alertSep()
	t.fmt.Fprintf(t.out, "EVENT: %s\n", v.Name)
	t.alertSep()
	t.renderNarrative(v.Narrative)
	t.thinSep()
	t.out.Flush()
}

func (t *Terminal) RenderEventResult(v view.EventResultView) {
	t.thinSep()
	t.renderNarrative(v.Narrative)
	t.thinSep()
	t.renderImpact("", v.Delta)
	t.ClearScreen()
}

func (t *Terminal) RenderInfo(msg string) {
	t.thinSep()
	t.out.WriteString(wrapText(msg, terminalWidth))
	t.linefeed()
	t.thinSep()
	t.waitForEnter()
	t.ClearScreen()
}

func (t *Terminal) RenderInfoNoWait(msg string) {
	t.thinSep()
	t.out.WriteString(wrapText(msg, terminalWidth))
	t.linefeed()
	t.thinSep()
	t.out.Flush()
}

func (t *Terminal) RenderEnding(v view.EndingView) {
	t.thickSep()
	t.renderNarrative(v.Narrative)
	t.fmt.Fprintf(t.out, "%s\n", wrapText("("+v.Explain+")", terminalWidth))
	t.thickSep()
	t.waitForEnterMsg("Game over. Press Enter to get back to main menu...")
	t.ClearScreen()
}

func (t *Terminal) renderNarrative(v []string) {
	l := len(v)
	if l == 0 {
		return
	}
	for _, line := range v[:l-1] {
		t.out.WriteString(wrapText(line, terminalWidth))
		t.linefeed()
		t.linefeed()
		t.waitForEnter()
		t.clearLine()
	}
	t.out.WriteString(wrapText(v[l-1], terminalWidth))
	t.linefeed()
	t.out.Flush()
}

func wrapText(text string, width int) string {
	if text == "" || width <= 0 {
		return text
	}

	var b strings.Builder
	lines := strings.Split(text, "\n")

	for li, line := range lines {
		words := strings.Fields(line)
		if len(words) == 0 {
			if li > 0 {
				b.WriteByte('\n')
			}
			continue
		}

		if li > 0 {
			b.WriteByte('\n')
		}

		lineLen := 0
		for i, word := range words {
			wordLen := len(word)

			if i == 0 {
				b.WriteString(word)
				lineLen = wordLen
				continue
			}

			if lineLen+1+wordLen > width {
				b.WriteByte('\n')
				b.WriteString(word)
				lineLen = wordLen
				continue
			}

			b.WriteByte(' ')
			b.WriteString(word)
			lineLen += 1 + wordLen
		}
	}

	return b.String()
}

func (t *Terminal) waitForEnterMsg(msg string) {
	t.out.WriteString(msg)
	t.out.Flush()
	t.in.ReadString('\n')
}

func (t *Terminal) waitForEnter() {
	t.waitForEnterMsg("(Press Enter to continue...)")
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

func (t *Terminal) alertSep() {
	t.out.Write(alertSep)
	t.linefeed()
}

func (t *Terminal) linefeed() {
	t.out.WriteByte('\n')
}
