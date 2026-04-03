package program

import (
	"reflect"
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
	"github.com/quangd42/silicon_valley_trail/internal/view"
	"github.com/quangd42/silicon_valley_trail/internal/weather"
)

type scriptedRenderer struct {
	tb               testing.TB
	prompts          []model.PromptChoice
	promptIndex      int
	calls            []string
	weathers         []model.WeatherKind
	actionResults    []view.ActionResultView
	eventInfos       []view.EventView
	eventResults     []view.EventResultView
	endings          []view.EndingView
	infos            []string
	clearScreenCount int
}

func (r *scriptedRenderer) RenderIntro(view.IntroView) {
	r.calls = append(r.calls, "RenderIntro")
}

func (r *scriptedRenderer) RenderDayInfo(v view.DayView) {
	r.calls = append(r.calls, "RenderDayInfo")
	r.weathers = append(r.weathers, v.Weather)
}

func (r *scriptedRenderer) RenderPrompt(view.PromptView) model.PromptChoice {
	r.calls = append(r.calls, "RenderPrompt")
	if r.promptIndex >= len(r.prompts) {
		r.tb.Helper()
		r.tb.Fatalf("unexpected prompt %d", r.promptIndex+1)
	}
	choice := r.prompts[r.promptIndex]
	r.promptIndex++
	return choice
}

func (r *scriptedRenderer) RenderActionResult(v view.ActionResultView) {
	r.calls = append(r.calls, "RenderActionResult")
	r.actionResults = append(r.actionResults, v)
}

func (r *scriptedRenderer) RenderEventInfo(v view.EventView) {
	r.calls = append(r.calls, "RenderEventInfo")
	r.eventInfos = append(r.eventInfos, v)
}

func (r *scriptedRenderer) RenderEventResult(v view.EventResultView) {
	r.calls = append(r.calls, "RenderEventResult")
	r.eventResults = append(r.eventResults, v)
}

func (r *scriptedRenderer) RenderInfo(msg string) {
	r.calls = append(r.calls, "RenderInfo")
	r.infos = append(r.infos, msg)
}

func (r *scriptedRenderer) RenderInfoNoWait(string) {
	r.calls = append(r.calls, "RenderInfoNoWait")
}

func (r *scriptedRenderer) RenderEnding(v view.EndingView) {
	r.calls = append(r.calls, "RenderEnding")
	r.endings = append(r.endings, v)
}

func (r *scriptedRenderer) ClearScreen() {
	r.calls = append(r.calls, "ClearScreen")
	r.clearScreenCount++
}

type seqRNG struct {
	rolls []int
	index int
}

func (r *seqRNG) IntN(max int) int {
	if r.index >= len(r.rolls) {
		panic("rng exhausted")
	}
	roll := r.rolls[r.index]
	r.index++
	if roll < 0 || roll >= max {
		panic("rng roll out of range")
	}
	return roll
}

type stubSaver struct {
	saveCount int
}

func (s *stubSaver) Save(*model.State) error {
	s.saveCount++
	return nil
}

func (s *stubSaver) Load(*model.State) error {
	return nil
}

func testDefinition() *gamedef.Definition {
	def := gamedef.Load()
	def.Route = []model.Location{
		{ID: "san-jose", Name: "San Jose"},
		{ID: "san-francisco", Name: "San Francisco"},
	}
	def.Events = map[string]gamedef.EventData{
		"helpful-founder": {
			ID:        "helpful-founder",
			Name:      "Helpful founder",
			Narrative: gamedef.Narrative{"A founder from the last accelerator batch spots you."},
			Choices: []gamedef.EventChoiceData{
				{
					Name:      "Take the bridge loan",
					Narrative: gamedef.Narrative{"You take the cash and keep moving."},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{Cash: 50},
						}
					},
				},
			},
		},
	}
	def.EventPools = model.EventPools{
		Main: []string{"helpful-founder"},
	}
	return def
}

func countCalls(calls []string, want string) int {
	count := 0
	for _, got := range calls {
		if got == want {
			count++
		}
	}
	return count
}

func Test_startGame(t *testing.T) {
	t.Run("weather refreshes every day", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(
			def.Route,
			// skipping event
			nil,
			nil,
		)
		state.Resources = model.Resources{
			Cash:   20000,
			Morale: 100,
			Coffee: 300,
		}

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionRest},
				{Kind: model.ChoiceAction, Action: model.ActionBuild},
				{Kind: model.ChoiceAction, Action: model.ActionMarket},
				{Kind: model.ChoiceControl, Control: model.ControlQuitToMenu},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      &seqRNG{},
			def:      def,
			weather:  weather.NewMockService(),
		}

		prog.startGame(state, true)

		wantCalls := []string{
			"ClearScreen",
			"RenderIntro",
			"RenderDayInfo",
			"RenderPrompt",
			"RenderActionResult",
			"RenderDayInfo",
			"RenderPrompt",
			"RenderActionResult",
			"RenderDayInfo",
			"RenderPrompt",
			"RenderActionResult",
			"RenderDayInfo",
			"RenderPrompt",
			"ClearScreen",
		}
		if !reflect.DeepEqual(renderer.calls, wantCalls) {
			t.Fatalf("call order = %#v, want %#v", renderer.calls, wantCalls)
		}
		wantWeathers := []model.WeatherKind{
			model.WeatherClear,
			model.WeatherRainy,
			model.WeatherFog,
			model.WeatherCloudy,
		}
		if !reflect.DeepEqual(renderer.weathers, wantWeathers) {
			t.Fatalf("weather renders = %#v, want %#v", renderer.weathers, wantWeathers)
		}
		if len(renderer.eventInfos) != 0 || len(renderer.eventResults) != 0 {
			t.Fatalf("event renders = (%d info, %d result), want (0, 0)", len(renderer.eventInfos), len(renderer.eventResults))
		}
	})

	t.Run("ending evaluates before arrival event", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(
			def.Route,
			// skipping event
			nil,
			nil,
		)
		state.Resources = model.Resources{
			Cash:   300,
			Morale: 100,
			Coffee: 30,
		}

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionTravel},
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		rng := &seqRNG{}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      rng,
			def:      def,
		}

		prog.startGame(state, false)

		wantCalls := []string{
			"ClearScreen",
			"RenderDayInfo",
			"RenderPrompt",
			"RenderActionResult",
			"RenderEnding",
		}
		if !reflect.DeepEqual(renderer.calls, wantCalls) {
			t.Fatalf("call order = %#v, want %#v", renderer.calls, wantCalls)
		}
		if len(renderer.eventInfos) != 0 || len(renderer.eventResults) != 0 {
			t.Fatalf("event renders = (%d info, %d result), want (0, 0)", len(renderer.eventInfos), len(renderer.eventResults))
		}
		if len(renderer.endings) != 1 {
			t.Fatalf("ending renders = %d, want 1", len(renderer.endings))
		}
		if got := renderer.endings[0].Explain; got != def.Endings[logic.EndingNoCash].Explain {
			t.Fatalf("ending explain = %q, want %q", got, def.Endings[logic.EndingNoCash].Explain)
		}
		if state.Resources.Cash != 0 {
			t.Fatalf("Cash = %d, want 0", state.Resources.Cash)
		}
	})

	t.Run("surviving travel plays arrival event before the final ending", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.Resources = model.Resources{
			Cash:    1000,
			Morale:  100,
			Coffee:  30,
			Product: 20,
			Hype:    10,
		}

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionTravel},
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		rng := &seqRNG{rolls: []int{
			0,  // select event
			99, // roll final pitch -> losing
		}}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      rng,
			def:      def,
		}
		prog.startGame(state, false)

		wantCalls := []string{
			"ClearScreen",
			"RenderDayInfo",
			"RenderPrompt",
			"RenderActionResult",
			"RenderEventInfo",
			"RenderPrompt",
			"RenderEventResult",
			"RenderEnding",
		}
		if !reflect.DeepEqual(renderer.calls, wantCalls) {
			t.Fatalf("call order = %#v, want %#v", renderer.calls, wantCalls)
		}
		if len(renderer.eventInfos) != 1 || len(renderer.eventResults) != 1 {
			t.Fatalf("event renders = (%d info, %d result), want (1, 1)", len(renderer.eventInfos), len(renderer.eventResults))
		}
		if len(renderer.endings) != 1 {
			t.Fatalf("ending renders = %d, want 1", len(renderer.endings))
		}
		if got := renderer.endings[0].Explain; got != def.Endings[logic.EndingNoOffer].Explain {
			t.Fatalf("ending explain = %q, want %q", got, def.Endings[logic.EndingNoOffer].Explain)
		}
		if state.Resources.Cash != 750 {
			t.Fatalf("Cash = %d, want 750", state.Resources.Cash)
		}
	})

	t.Run("non-travel action does not trigger arrival event", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(
			def.Route,
			// event not triggered
			nil,
			nil,
		)

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionRest},
				{Kind: model.ChoiceControl, Control: model.ControlQuitToMenu},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      &seqRNG{},
			def:      def,
		}

		prog.startGame(state, false)

		if len(renderer.eventInfos) != 0 || len(renderer.eventResults) != 0 {
			t.Fatalf("unexpected event renders: (%d info, %d result)", len(renderer.eventInfos), len(renderer.eventResults))
		}
		if len(renderer.endings) != 0 {
			t.Fatalf("ending renders = %d, want 0", len(renderer.endings))
		}
		if renderer.clearScreenCount != 2 {
			t.Fatalf("ClearScreen count = %d, want 2", renderer.clearScreenCount)
		}
	})

	t.Run("saved game made during event will start at event", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.CurrentLocation = 1
		state.CurrentEvent = "helpful-founder"

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		saver := &stubSaver{}
		prog := &Program{
			renderer: renderer,
			saver:    saver,
			rng: &seqRNG{rolls: []int{
				// IMPORTANT that we only need the roll for the final pitch here, not
				// for event roll because event is saved and restored
				99, // roll final pitch -> losing
			}},
			def: def,
		}

		prog.startGame(state, false)

		if len(renderer.eventInfos) != 1 {
			t.Fatalf("event info renders = %d, want 1", len(renderer.eventInfos))
		}
		if len(renderer.eventResults) != 1 {
			t.Fatalf("event result renders = %d, want 1", len(renderer.eventResults))
		}
		if state.Resources.Cash != 7050 {
			t.Fatalf("Cash = %d, want 7050", state.Resources.Cash)
		}
		if len(renderer.endings) != 1 {
			t.Fatalf("ending renders = %d, want 1", len(renderer.endings))
		}
		if got := renderer.endings[0].Explain; got != def.Endings[logic.EndingNoOffer].Explain {
			t.Fatalf("ending explain = %q, want %q", got, def.Endings[logic.EndingNoOffer].Explain)
		}
		if state.CurrentLocation != len(state.Route)-1 {
			t.Fatalf("CurrentLocation = %d, want %d", state.CurrentLocation, len(state.Route)-1)
		}
	})

	t.Run("play event skips gracefully when there is no more events to play", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(
			def.Route,
			// event not triggered
			nil,
			nil,
		)

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionTravel},
				{Kind: model.ChoiceControl, Control: model.ControlQuitToMenu},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng: &seqRNG{rolls: []int{
				0, // Roll final pitch roll -> winning
			}},
			def: def,
		}

		prog.startGame(state, false)

		if len(renderer.eventInfos) != 0 {
			t.Fatalf("event info renders = %d, want 0", len(renderer.eventInfos))
		}
		if len(renderer.eventResults) != 0 {
			t.Fatalf("event result renders = %d, want 0", len(renderer.eventResults))
		}
		if len(renderer.endings) != 1 {
			t.Fatalf("ending renders = %d, want 1", len(renderer.endings))
		}
		if got := renderer.endings[0].Explain; got != def.Endings[logic.EndingTogether].Explain {
			t.Fatalf("ending explain = %q, want %q", got, def.Endings[logic.EndingTogether].Explain)
		}
		if state.CurrentLocation != len(state.Route)-1 {
			t.Fatalf("CurrentLocation = %d, want %d", state.CurrentLocation, len(state.Route)-1)
		}
	})
	t.Run("reaching the destination resolves the final ending path", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.Resources = model.Resources{
			Cash:    1000,
			Morale:  100,
			Coffee:  30,
			Product: 20,
			Hype:    10,
		}

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionTravel},
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		rng := &seqRNG{rolls: []int{
			0, // select event
			0, // roll final pitch -> winning
		}}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      rng,
			def:      def,
		}

		prog.startGame(state, false)

		if len(renderer.eventInfos) != 1 || len(renderer.eventResults) != 1 {
			t.Fatalf("unexpected event renders: (%d info, %d result)", len(renderer.eventInfos), len(renderer.eventResults))
		}
		if len(renderer.endings) != 1 {
			t.Fatalf("ending renders = %d, want 1", len(renderer.endings))
		}
		if got := renderer.endings[0].Explain; got != def.Endings[logic.EndingTogether].Explain {
			t.Fatalf("ending explain = %q, want %q", got, def.Endings[logic.EndingTogether].Explain)
		}
		if state.CurrentLocation != len(state.Route)-1 {
			t.Fatalf("CurrentLocation = %d, want %d", state.CurrentLocation, len(state.Route)-1)
		}
	})

	t.Run("new game seeds events from the injected definition", func(t *testing.T) {
		def := testDefinition()

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceAction, Action: model.ActionTravel},
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng: &seqRNG{rolls: []int{
				0,  // select event
				99, // roll final pitch -> losing
			}},
			def: def,
		}

		prog.newGame()

		if len(renderer.eventInfos) != 1 {
			t.Fatalf("event info renders = %d, want 1", len(renderer.eventInfos))
		}
		if got := renderer.eventInfos[0].Name; got != "Helpful founder" {
			t.Fatalf("event name = %q, want %q", got, "Helpful founder")
		}
		if len(renderer.eventResults) != 1 {
			t.Fatalf("event result renders = %d, want 1", len(renderer.eventResults))
		}
		if got := renderer.eventResults[0].Delta.Cash; got != 50 {
			t.Fatalf("event result cash delta = %d, want 50", got)
		}
	})
}

func Test_playEvent(t *testing.T) {
	t.Run("in WeatherFog, roll above event threshold skips the event entirely", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(
			def.Route,
			// skipping event
			nil,
			nil,
		)
		state.CurrentLocation = 1
		state.Weather = model.WeatherFog

		renderer := &scriptedRenderer{tb: t}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      &seqRNG{rolls: []int{1}}, // Roll for WeatherFog
			def:      def,
		}

		got := prog.playEvent(state)

		if got {
			t.Fatal("playEvent() = true, want false")
		}
		if len(renderer.calls) != 0 {
			t.Fatalf("unexpected renderer calls = %#v", renderer.calls)
		}
	})

	t.Run("event is taken out of the pool after being selected", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.CurrentLocation = 1
		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      &seqRNG{rolls: []int{0}}, // Roll to select event
			def:      def,
		}

		got := prog.playEvent(state)

		if got {
			t.Fatal("playEvent() = true, want false")
		}
		if len(renderer.eventInfos) != 1 {
			t.Fatalf("unexpected renderer calls = %#v", renderer.calls)
		}
		if len(renderer.eventResults) != 1 {
			t.Fatalf("event result renders = %d, want 1", len(renderer.eventResults))
		}
		if len(state.EventPools.Main) != 0 {
			t.Fatalf("event pool count = %d, want 0", len(state.EventPools.Main))
		}
	})

	t.Run("save keeps the player inside the event loop", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.CurrentLocation = 1

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceControl, Control: model.ControlSave},
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		saver := &stubSaver{}
		prog := &Program{
			renderer: renderer,
			saver:    saver,
			rng: &seqRNG{rolls: []int{
				0, // Roll to select event
			}},
			def: def,
		}

		got := prog.playEvent(state)

		if got {
			t.Fatal("playEvent() = true, want false")
		}
		if saver.saveCount != 1 {
			t.Fatalf("Save count = %d, want 1", saver.saveCount)
		}
		if len(renderer.eventInfos) != 2 {
			t.Fatalf("event info renders = %d, want 2", len(renderer.eventInfos))
		}
		if len(renderer.eventResults) != 1 {
			t.Fatalf("event result renders = %d, want 1", len(renderer.eventResults))
		}
		if len(renderer.infos) != 1 || renderer.infos[0] != "Game saved." {
			t.Fatalf("info messages = %#v, want [\"Game saved.\"]", renderer.infos)
		}
		if state.Resources.Cash != 7050 {
			t.Fatalf("Cash = %d, want 7050", state.Resources.Cash)
		}
	})

	t.Run("quit to menu exits the event loop immediately", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.CurrentLocation = 1

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceControl, Control: model.ControlQuitToMenu},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng: &seqRNG{rolls: []int{
				0, // Roll to select event
			}},
			def: def,
		}

		got := prog.playEvent(state)

		if !got {
			t.Fatal("playEvent() = false, want true")
		}
		if len(renderer.eventInfos) != 1 {
			t.Fatalf("event info renders = %d, want 1", len(renderer.eventInfos))
		}
		if len(renderer.eventResults) != 0 {
			t.Fatalf("event result renders = %d, want 0", len(renderer.eventResults))
		}
		if renderer.clearScreenCount != 1 {
			t.Fatalf("ClearScreen count = %d, want 1", renderer.clearScreenCount)
		}
	})

	t.Run("weather-conditioned events are selected before the main pool", func(t *testing.T) {
		def := testDefinition()
		def.Events["clear-skies"] = gamedef.EventData{
			ID:        "clear-skies",
			Name:      "Clear Skies",
			Narrative: gamedef.Narrative{"The sun is doing free marketing for you."},
			Choices: []gamedef.EventChoiceData{
				{
					Name:      "Take the boost",
					Narrative: gamedef.Narrative{"You lean into the moment."},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{Hype: 3},
						}
					},
				},
			},
			Conditions: gamedef.EventConditions{
				Weather: model.WeatherClear,
			},
		}
		def.EventPools = model.EventPools{
			Main: []string{"helpful-founder"},
			Weather: map[model.WeatherKind][]string{
				model.WeatherClear: {"clear-skies"},
			},
		}
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.CurrentLocation = 1
		state.Weather = model.WeatherClear

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      &seqRNG{rolls: []int{0}},
			def:      def,
		}

		got := prog.playEvent(state)

		if got {
			t.Fatal("playEvent() = true, want false")
		}
		if len(renderer.eventInfos) != 1 {
			t.Fatalf("event info renders = %d, want 1", len(renderer.eventInfos))
		}
		if got := renderer.eventInfos[0].Name; got != "Clear Skies" {
			t.Fatalf("event name = %q, want %q", got, "Clear Skies")
		}
		if len(state.EventPools.Main) != 1 || state.EventPools.Main[0] != "helpful-founder" {
			t.Fatalf("main pool = %#v, want %#v", state.EventPools.Main, []string{"helpful-founder"})
		}
		if len(state.EventPools.Weather[model.WeatherClear]) != 0 {
			t.Fatalf("weather pool count = %d, want 0", len(state.EventPools.Weather[model.WeatherClear]))
		}
	})

	t.Run("falls back to the main pool when the current weather pool is empty", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)
		state.CurrentLocation = 1
		state.Weather = model.WeatherClear

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceEvent, EventChoiceIndex: 0},
			},
		}
		prog := &Program{
			renderer: renderer,
			saver:    &stubSaver{},
			rng:      &seqRNG{rolls: []int{0}},
			def:      def,
		}

		got := prog.playEvent(state)

		if got {
			t.Fatal("playEvent() = true, want false")
		}
		if len(renderer.eventInfos) != 1 {
			t.Fatalf("event info renders = %d, want 1", len(renderer.eventInfos))
		}
		if got := renderer.eventInfos[0].Name; got != "Helpful founder" {
			t.Fatalf("event name = %q, want %q", got, "Helpful founder")
		}
		if len(state.EventPools.Main) != 0 {
			t.Fatalf("main pool count = %d, want 0", len(state.EventPools.Main))
		}
		if len(state.EventPools.Weather[model.WeatherClear]) != 0 {
			t.Fatalf("weather pool count = %d, want 0", len(state.EventPools.Weather[model.WeatherClear]))
		}
	})
}

func Test_playTurn(t *testing.T) {
	t.Run("save loops back to the day prompt", func(t *testing.T) {
		def := testDefinition()
		state := model.NewState(def.Route, def.EventPools.Main, def.EventPools.Weather)

		renderer := &scriptedRenderer{
			tb: t,
			prompts: []model.PromptChoice{
				{Kind: model.ChoiceControl, Control: model.ControlSave},
				{Kind: model.ChoiceAction, Action: model.ActionRest},
			},
		}
		saver := &stubSaver{}
		prog := &Program{
			renderer: renderer,
			saver:    saver,
			rng:      &seqRNG{},
			def:      def,
		}

		got := prog.playTurn(state)

		if got.quitToMenu {
			t.Fatal("quitToMenu = true, want false")
		}
		if got.traveled {
			t.Fatal("traveled = true, want false")
		}
		if saver.saveCount != 1 {
			t.Fatalf("Save count = %d, want 1", saver.saveCount)
		}
		if countCalls(renderer.calls, "RenderDayInfo") != 2 {
			t.Fatalf("RenderDayInfo count = %d, want 2", countCalls(renderer.calls, "RenderDayInfo"))
		}
		if countCalls(renderer.calls, "RenderPrompt") != 2 {
			t.Fatalf("RenderPrompt count = %d, want 2", countCalls(renderer.calls, "RenderPrompt"))
		}
		if countCalls(renderer.calls, "RenderInfo") != 1 {
			t.Fatalf("RenderInfo count = %d, want 1", countCalls(renderer.calls, "RenderInfo"))
		}
		if countCalls(renderer.calls, "RenderActionResult") != 1 {
			t.Fatalf("RenderActionResult count = %d, want 1", countCalls(renderer.calls, "RenderActionResult"))
		}
		if len(renderer.infos) != 1 || renderer.infos[0] != "Game saved." {
			t.Fatalf("info messages = %#v, want [\"Game saved.\"]", renderer.infos)
		}
	})
}
