package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
)

type IntroView gamedef.Narrative

func Intro(def *gamedef.Definition) IntroView {
	return IntroView(def.Intro)
}

type EndingView gamedef.EndingCopy

func Ending(ending logic.Ending, def *gamedef.Definition) EndingView {
	return EndingView(def.Endings[ending])
}
