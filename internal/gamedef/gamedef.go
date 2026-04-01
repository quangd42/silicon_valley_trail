// Package gamedef defines the authored game data and scripted effects.
package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type Definition struct {
	Intro   []string
	Route   []model.Location
	Actions map[model.Action]ActionData
	Weather map[model.WeatherKind]WeatherData
	Endings map[logic.Ending]EndingCopy
}

// Load returns the authored game definition.
func Load() *Definition {
	return &Definition{
		Intro:   introCopy(),
		Route:   DefaultRoute(),
		Actions: actionData(),
		Weather: weatherCopy(),
		Endings: endingCopy(),
	}
}

func introCopy() []string {
	return []string{
		`Welcome to Silicon Valley Trail!

You and your best bud Pete set out from your HQ in San Jose to San Francisco to attend a
major investor meeting. Your product: a sleeping mask that lets people relive childhood
memories through dreams.

Will you be able to impress the investors?`,
		`Manage your resources wisely:
* Cash    ($)   : Don’t run out. No cash = game over.
* Morale  (%)   : Keep the team motivated.
* Coffee  (cups): Your startup fuel. 2 days without it = game over.
* Product (%)   : How ready your product is. Directly affects your odds of getting signed.
* Hype    (%)   : Public interest in your startup. Every 2 Hype = 1 Product.`,
	}
}

type Narrative []string

type EndingCopy struct {
	Narrative Narrative
	Desc      string
}

func endingCopy() map[logic.Ending]EndingCopy {
	return map[logic.Ending]EndingCopy{
		logic.EndingNone: {}, // This is just a placeholder, this is not a real ending
		logic.EndingNoCash: {
			Narrative: Narrative{"You have just enough to catch a train home. Time for some part-time jobs..."},
			Desc:      "You have no cash left.",
		},
		logic.EndingNoCoffee: {
			Narrative: Narrative{"*Confused Psyduck*."},
			Desc:      "You have no coffee for far too long.",
		},
		logic.EndingNoOffer: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco. After one last rushed coffee,\nyou step into the meeting room and face the investors.",
				"The pitch is solid, but the product is not. The investors don't seem convinced.",
				"...",
				"Maybe you got here too early.",
			},
			Desc: "The investors turned you down.",
		},
		logic.EndingAlone: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco. After one last rushed coffee,\nyou step into the meeting room and face the investors.",
				"The presentation lands. Against the odds, you leave with a verbal commitment.",
				"...",
				"But what did it cost?",
			},
			Desc: "You won over the investors.",
		},
		logic.EndingTogether: {}, // To be added if we ever get there
	}
}
