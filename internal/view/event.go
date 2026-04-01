package view

import "github.com/quangd42/silicon_valley_trail/internal/gamedef"

type EventView = gamedef.EventData

func Event(def *gamedef.Definition, eventID int) EventView {
	return def.Events[eventID]
}
