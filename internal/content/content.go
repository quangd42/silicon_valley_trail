// Package content loads the game copy into a structured copy pool to use across different parts of the game.
// It enables game copy to be loaded from a database or copy files later on. For now, it is embedded directly
// in the code.
package content

import "github.com/quangd42/silicon_valley_trail/internal/model"

type Content struct {
	Intro   string
	Route   []model.Location
	Actions map[model.Action]ActionCopy
	Ending  string
}

func Load() *Content {
	return &Content{
		Intro:   introCopy(),
		Route:   DefaultRoute(),
		Actions: actions(),
		Ending:  endingCopy(),
	}
}

func introCopy() string {
	return `Welcome to Silicon Valley Trail!

You and your best bud Pete set out from your HQ in San Jose to San Francisco to attend a major investor meeting.
Your product: a sleeping mask that lets people relive childhood memories through dreams.

Will you be able to impress the investors?

`
}

func endingCopy() string {
	return `Congratulations! You reached San Francisco! After a hasty coffee you went into the meeting room where investors are waiting for you.

The presentation went well. You got a verbal contract on the spot... Your team is estatic!

...

But what did it cost?...`
}

type ActionCopy struct {
	Desc   string
	Result string
}

func actions() map[model.Action]ActionCopy {
	return map[model.Action]ActionCopy{
		model.ActionTravel: {
			Desc:   "Travel to the next location",
			Result: "Your team hit the road...\n\nArrived at %s!\n\n",
		},
		model.ActionRest: {
			Desc:   "Rest and recover (restore morale, use coffee)",
			Result: "You decided to take a break...\n\nYou're rested, and filled with determination.\n\n",
		},
		model.ActionBuild: {
			Desc:   "Work on product (reduce bugs, use coffee)",
			Result: "You take on the next item on the roadmap...\n\nYou're happy with the result, but everyone is tired...\n(Product is %d%% more ready. Cost %d coffee, %d%% morale)\n\n",
		},
		model.ActionMarket: {
			Desc:   "Marketing push (increase hype, costs money)",
			Result: "You launch a marketing campaign...\n\nEvery \"debate\" on X is about your product...\n(Hype increased. Cost $%d)\n\n",
		},
		model.ActionSave: {
			Desc: "Save game",
		},
		model.ActionQuit: {
			Desc: "Quit to menu",
		},
	}
}
