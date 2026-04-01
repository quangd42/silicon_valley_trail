package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type EventData struct {
	Name         string
	Narrative    Narrative
	ChoicesLabel string
	Choices      []EventChoiceData
}

type EventChoiceData struct {
	Name      string
	Desc      string
	Narrative Narrative
	Effect    logic.Effect
}

func eventData() []EventData {
	return []EventData{
		{
			Name: "We meet again!",
			Narrative: Narrative{
				"A cheery disheveled fellow approaches you gleefully. You do not know this man.",
				`"It's me, Ranwid! Have any goods for me today? The usual? A fella like me can't make it alone, you know?`,
				"You eye him suspiciously and consider your options...",
			},
			ChoicesLabel: "Give...",
			Choices: []EventChoiceData{
				{
					Name: "...some coffee.",
					Desc: "-2 Coffee, +10 Hype",
					Narrative: Narrative{
						`Ranwid: "Exquisite! Was feeling parched."

*Glup glup glup*

He downs the coffee in one go and lets out a satisfied burp.`,
						`He rummages around his various pockets...

Ranwid: "Here, look what I've got for you today! Take it take it!"`,
					},
					Effect: func(s *model.State, ctx logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +10,
							},
						}
					},
				},
				{
					Name: "...some money.",
					Desc: "-200 Cash, +10 Hype",
					Narrative: Narrative{
						`Ranwid: "Magnificent! This will be quite handy if I run into those *mask wearing hoodlums* again."`,
						`He rummages around his various pockets...

Ranwid: "Here, look what I've got for you today! Take it take it!"`,
					},
					Effect: func(s *model.State, ctx logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -200,
								Hype: +10,
							},
						}
					},
				},
				{
					Name: "Ignore him",
					Desc: "...",
					Narrative: Narrative{
						`Ranwid: "Aaaaagghh!! What a jerk you are sometimes!"

He runs away.`,
					},
					Effect: func(s *model.State, ctx logic.Context) logic.Change {
						return logic.Change{}
					},
				},
			},
		},
	}
}
