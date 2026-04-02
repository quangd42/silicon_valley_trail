package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type WeatherData struct {
	Desc   string
	Effect logic.Effect
}

func weatherCopy() map[model.WeatherKind]WeatherData {
	return map[model.WeatherKind]WeatherData{
		model.WeatherUnknown: {
			Desc: "What is going on?",
		},
		model.WeatherClear: {
			Desc: "You feel productive and ready to go.\n(1.Travel+  3.Build+)",
			Effect: func(s *model.State, ctx logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				switch ctx.Action {
				case model.ActionTravel:
					return logic.Change{
						Delta: model.Resources{Morale: 2},
					}
				case model.ActionBuild:
					return logic.Change{
						Delta: model.Resources{
							Morale:  3,
							Product: teamSize * s.Resources.Morale / 100,
						},
					}
				default:
					return logic.Change{}
				}
			},
		},
		model.WeatherRainy: {
			Desc: "It's miserable out there.\n(1.Travel-  4.Marketing-)",
			Effect: func(s *model.State, ctx logic.Context) logic.Change {
				teamSize := len(s.Party.Members)
				switch ctx.Action {
				case model.ActionTravel:
					return logic.Change{
						Delta: model.Resources{
							Coffee: -teamSize,
							Morale: -4,
						},
					}
				case model.ActionMarket:
					return logic.Change{
						Delta: model.Resources{Hype: -3},
					}
				default:
					return logic.Change{}
				}
			},
		},
		model.WeatherFog: {
			Desc: "You feel unsure about the future.\n(1.Travel-  3.Build-)",
			Effect: func(s *model.State, ctx logic.Context) logic.Change {
				switch ctx.Action {
				case model.ActionTravel:
					return logic.Change{
						Delta: model.Resources{Morale: -3},
					}
				case model.ActionBuild:
					return logic.Change{
						Delta: model.Resources{Product: -2},
					}
				default:
					return logic.Change{}
				}
			},
		},
		model.WeatherCloudy: {
			Desc: "Everything feels in place.\n(2.Rest+  4.Marketing+)",
			Effect: func(s *model.State, ctx logic.Context) logic.Change {
				switch ctx.Action {
				case model.ActionRest:
					return logic.Change{
						Delta: model.Resources{Morale: 4, Coffee: 2},
					}
				case model.ActionMarket:
					return logic.Change{
						Delta: model.Resources{Hype: 3},
					}
				default:
					return logic.Change{}
				}
			},
		},
	}
}
