// Package gamedef defines the authored game data and scripted effects.
package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type Definition struct {
	Intro       Narrative
	Route       []model.Location
	Actions     map[model.Action]ActionData
	ActionOrder []model.Action
	Weather     map[model.WeatherKind]WeatherData
	Events      map[string]EventData
	EventIDs    []string
	Endings     map[logic.Ending]EndingCopy
}

type Narrative []string

// Load returns the authored game definition.
func Load() *Definition {
	return &Definition{
		Intro:       introCopy(),
		Route:       DefaultRoute(),
		Actions:     actionData(),
		ActionOrder: actionOrder(),
		Weather:     weatherCopy(),
		Events:      eventData(),
		EventIDs:    eventIDs(),
		Endings:     endingCopy(),
	}
}

func introCopy() Narrative {
	return Narrative{
		`Welcome to Silicon Valley Trail!

		You and your best bud Pete set out from your HQ in San Jose to San Francisco to attend a major investor meeting. Your product: a sleeping mask that lets people relive childhood memories through dreams.

		Will you be able to impress the investors?`,
		"Instructions:\n\nArrive at San Francisco with the best possible product - the better it is, the higher the chance you'll win over the investors in San Francisco.",
		`Manage your resources wisely:
		* Cash    ($): Don’t run out. No cash = game over.
		* Morale  (%): Motivated team improve product faster.
		* Coffee  (ct): Your startup fuel. 2 days without it = game over.
		* Product (%): How ready your product is = your odds of getting signed.
		* Hype    (%): Public interest in your startup. Every 2 Hype = 1 Product.`,
	}
}

type EndingCopy struct {
	Narrative Narrative
	Explain   string
}

func endingCopy() map[logic.Ending]EndingCopy {
	return map[logic.Ending]EndingCopy{
		logic.EndingNone: {}, // This is just a placeholder, this is not a real ending
		logic.EndingNoCash: {
			Narrative: Narrative{"You have just enough to catch a train home. Time for some part-time jobs..."},
			Explain:   "You have no cash left.",
		},
		logic.EndingNoCoffee: {
			Narrative: Narrative{"*Confused Psyduck*."},
			Explain:   "You have no coffee for far too long.",
		},
		logic.EndingNoOffer: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco. After one last rushed coffee, you step into the meeting room and face the investors.",
				"The pitch is solid, but the product is not. The investors don't seem convinced.",
				"...",
				"Maybe you got here too early.",
			},
			Explain: "The investors turned you down.",
		},
		logic.EndingTogether: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco. After one last rushed coffee, you step into the meeting room and face the investors.",
				"The presentation lands. Against the odds, you leave with a verbal commitment.",
				"...",
				"But what did it cost?",
			},
			Explain: "You won over the investors.",
		},
		logic.EndingAlone: {}, // To be added if we ever get there
	}
}
