// Package content loads the game copy into a structured copy pool to use across different parts of the game.
// It enables game copy to be loaded from a database or copy files later on. For now, it is embedded directly
// in the code.
package content

import "github.com/quangd42/silicon_valley_trail/internal/model"

type Content struct {
	Intro   []string
	Route   []model.Location
	Actions map[model.Action]ActionCopy
	Ending  string
}

// Load returns the structured copy pool for the entire game.
// Later on it can be improved to orchestrate fetching content from different sources.
func Load() *Content {
	return &Content{
		Intro:   introCopy(),
		Route:   DefaultRoute(),
		Actions: actions(),
		Ending:  endingCopy(),
	}
}

func introCopy() []string {
	return []string{
		`Welcome to Silicon Valley Trail!

You and your best bud Pete set out from your HQ in San Jose to San Francisco to attend a
major investor meeting. Your product: a sleeping mask that lets people relive childhood
memories through dreams.

Will you be able to impress the investors?
`,
		`
Manage your resources wisely:
* Cash    ($)   : Don’t run out. No cash = game over.
* Morale  (%)   : Keep the team motivated.
* Coffee  (cups): Your startup fuel. 2 days without it = game over.
* Product (%)   : How ready your product is. Directly affects your odds of getting signed.
* Hype    (%)   : Public interest in your startup. Every 2 Hype = 1 Product.
`,
	}
}

func endingCopy() string {
	return `Congratulations! You reached San Francisco! After a hasty coffee you went into the meeting room where investors are waiting for you.

The presentation went well. You got a verbal contract on the spot... Your team is estatic!

...

But what did it cost?...`
}

type ActionCopy struct {
	Desc      string
	Narrative []string
}

// TODO: when `actions` become interface, need unit tests to make sure every
func actions() map[model.Action]ActionCopy {
	return map[model.Action]ActionCopy{
		model.ActionTravel: {
			Desc:      "Travel to the next location",
			Narrative: []string{"Your team hit the road...\n\n"},
		},
		model.ActionRest: {
			Desc:      "Rest and recover (restore morale, use coffee)",
			Narrative: []string{"You decided to take a break...\n\n", "...zZz...\n\n", "You feel refreshed. You're filled with determination.\n\n"},
		},
		model.ActionBuild: {
			Desc:      "Work on product (reduce bugs, use coffee)",
			Narrative: []string{"You take on the next item on the roadmap...\n\n", "You're happy with the result, but everyone is tired...\n\n"},
		},
		model.ActionMarket: {
			Desc:      "Marketing push (increase hype, costs money)",
			Narrative: []string{"You launch a marketing campaign...\n\n", "Every \"debate\" on X is about your product.\n\n"},
		},
	}
}
