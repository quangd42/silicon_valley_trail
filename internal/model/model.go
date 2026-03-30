package model

type State struct {
	Day             int         `json:"day"`
	Route           []Location  `json:"route"`
	CurrentLocation int         `json:"current_location"`
	Resources       Resources   `json:"resources"`
	Party           Party       `json:"party"`
	Weather         WeatherKind `json:"weather"`
}

func NewState(route []Location) *State {
	return &State{
		Day:       0,
		Route:     route,
		Resources: defaultResources(),
		Party:     defaultParty(),
		Weather:   WeatherClear,
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
	Cash      int `json:"cash"`      // dollar
	Morale    int `json:"morale"`    // percent
	Coffee    int `json:"coffee"`    // cup
	Hype      int `json:"hype"`      // percent
	Readiness int `json:"readiness"` // percent
}

func defaultResources() Resources {
	return Resources{
		Cash:      10_000,
		Morale:    100,
		Coffee:    30,
		Hype:      10,
		Readiness: 20,
	}
}

func (r *Resources) Add(delta Resources) {
	r.Cash = addClamped(r.Cash, delta.Cash)
	r.Morale = addClamped(r.Morale, delta.Morale)
	r.Coffee = addClamped(r.Coffee, delta.Coffee)
	r.Hype = addClamped(r.Hype, delta.Hype)
	r.Readiness = addClamped(r.Readiness, delta.Readiness)
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

type WeatherKind string

const (
	WeatherClear  WeatherKind = "clear"
	WeatherRainy  WeatherKind = "rain"
	WeatherFog    WeatherKind = "fog"
	WeatherHot    WeatherKind = "heat"
	WeatherCloudy WeatherKind = "cloudy"
)

type Action int

const (
	ActionTravel Action = iota
	ActionRest
	ActionBuild
	ActionMarket

	// Sentinel value, holds the count of in-game actions. Mainly useful for testing.
	ActionCount
)

type Control string

const (
	ControlSave       Control = "Save Game"
	ControlLoad       Control = "Load Game"
	ControlQuitToMenu Control = "Quit to Menu"
	ControlNewGame    Control = "New Game"
	ControlQuitGame   Control = "Quit Game"
)
