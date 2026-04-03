package model

import "slices"

type State struct {
	Day              int         `json:"day"`
	Route            []Location  `json:"route"`
	CurrentLocation  int         `json:"current_location"`
	Resources        Resources   `json:"resources"`
	Party            Party       `json:"party"`
	Weather          WeatherKind `json:"weather"`
	NoCoffeeDayCount int         `json:"no_coffee_day_count"`
	EventPools       EventPools  `json:"event_pools"`
	// Indicates which event is in play.
	// An empty string indicates no event is in play.
	CurrentEvent string `json:"current_event"`
}

func NewState(
	route []Location,
	mainPool []string,
	weatherPools map[WeatherKind][]string,
) *State {
	pools := cloneEventPools(EventPools{
		Main:    mainPool,
		Weather: weatherPools,
	})
	return &State{
		Day:        1,
		Route:      route,
		Resources:  defaultResources(),
		Party:      defaultParty(),
		Weather:    WeatherClear,
		EventPools: pools,
	}
}

type Location struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Resources struct {
	Cash    int `json:"cash"`    // dollar
	Morale  int `json:"morale"`  // percent
	Coffee  int `json:"coffee"`  // cup
	Hype    int `json:"hype"`    // percent
	Product int `json:"product"` // percent
}

func defaultResources() Resources {
	return Resources{
		Cash:    6_000,
		Morale:  100,
		Coffee:  26,
		Hype:    10,
		Product: 20,
	}
}

func (r *Resources) Add(delta Resources) {
	r.Cash += delta.Cash
	r.Morale += delta.Morale
	r.Coffee += delta.Coffee
	r.Hype += delta.Hype
	r.Product += delta.Product
}

func (r *Resources) AddClamped(delta Resources) {
	r.Cash = addClamped(r.Cash, delta.Cash)
	r.Morale = addClamped(r.Morale, delta.Morale)
	r.Coffee = addClamped(r.Coffee, delta.Coffee)
	r.Hype = addClamped(r.Hype, delta.Hype)
	r.Product = addClamped(r.Product, delta.Product)
}

func addClamped(a, b int) int {
	sum := a + b
	if sum < 0 {
		return 0
	}
	return sum
}

type Party struct {
	Members []PartyMember `json:"members"`
}

func defaultParty() Party {
	return Party{
		Members: []PartyMember{
			{"You"},
			{"Pete"},
		},
	}
}

type PartyMember struct {
	Name string `json:"name"`
}

type WeatherKind int

const (
	WeatherUnknown WeatherKind = iota
	WeatherClear
	WeatherRainy
	WeatherFog
	WeatherCloudy

	// Sentinel value, holds the count of weather kinds. Mainly useful for testing.
	WeatherKindCount
)

func (w WeatherKind) String() string {
	switch w {
	case WeatherUnknown:
		return "Unknown"
	case WeatherClear:
		return "Clear"
	case WeatherRainy:
		return "Rainy"
	case WeatherFog:
		return "Fog"
	case WeatherCloudy:
		return "Cloudy"
	default:
		panic("unreachable")
	}
}

type Action int

const (
	ActionTravel Action = iota
	ActionRest
	ActionBuild
	ActionMarket

	// Sentinel value, holds the count of in-game actions. Mainly useful for testing.
	ActionCount
)

type Control int

const (
	ControlSave Control = iota
	ControlLoad
	ControlQuitToMenu
	ControlNewGame
	ControlQuitGame
)

func (c Control) String() string {
	switch c {
	case ControlSave:
		return "Save Game"
	case ControlLoad:
		return "Load Game"
	case ControlQuitToMenu:
		return "Quit to Menu"
	case ControlNewGame:
		return "New Game"
	case ControlQuitGame:
		return "Quit Game"
	default:
		panic("unreachable")
	}
}

// PromptChoice is the type of the data returned from prompting the player. It is a discriminated
// struct, to distinguish if the user has chosen an in-game action, event choice or a game session
// control. Even though it is defined here, its existence comes from the fact that we unfortunately
// have to mix those choices in the same input interface in the CLI UI representation, and it will
// need to be refactored if we ever add another UI renderer.
//
// When `Kind` = `ChoiceEvent`, `PromptChoice.EventChoiceIndex` is set to the index of the
// selected `EventChoice` in the current event.
// Accessing the unset field does not panic, simply returns the default (and wrong) value.
type PromptChoice struct {
	Kind             ChoiceKind
	Action           Action
	Control          Control
	EventChoiceIndex int
}

type ChoiceKind int

const (
	ChoiceControl ChoiceKind = iota
	ChoiceAction
	ChoiceEvent
)

type EventPools struct {
	Main    []string                 `json:"main"`
	Weather map[WeatherKind][]string `json:"weather"`
}

// cloneEventPools deep-copies event pool slices so each game state owns its pools.
func cloneEventPools(src EventPools) EventPools {
	out := EventPools{
		Main: slices.Clone(src.Main),
	}
	if src.Weather == nil {
		return out
	}
	out.Weather = make(map[WeatherKind][]string, len(src.Weather))
	for weather, pool := range src.Weather {
		out.Weather[weather] = slices.Clone(pool)
	}
	return out
}
