package model

type State struct {
	Day             int
	Route           []Location
	CurrentLocation int
	Resources       Resources
	Party           Party
	Weather         WeatherKind
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
	ID        string
	Name      string
	Desc      string
	Latitude  float64
	Longitude float64
}

type Resources struct {
	Cash      int
	Morale    int
	Coffee    int
	Hype      int
	Readiness int
}

func defaultResources() Resources {
	return Resources{
		Cash:      30_000,
		Morale:    100,
		Coffee:    30,
		Hype:      30,
		Readiness: 20,
	}
}

type Party struct {
	Members []PartyMember
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
	Name string
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
	ActionSave
	ActionQuit
)
