package model

type ActionKind string

// The copy of each action could be stored in `content` package. Hardcoded
// for now for simplicity.
const (
	ActionTravel ActionKind = "Travel to the next location"
	ActionRest   ActionKind = "Rest and recover (restore morale, use coffee)"
	ActionBuild  ActionKind = "Work on product (reduce bugs, use coffee)"
	ActionMarket ActionKind = "Marketing push (increase hype, costs money)"
	ActionSave   ActionKind = "Save game"
	ActionQuit   ActionKind = "Quit to menu"
)
