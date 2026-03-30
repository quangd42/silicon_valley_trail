package view

import (
	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/logic"
)

type IntroView content.Narrative

func Intro(cont *content.Content) IntroView {
	return IntroView(cont.Intro)
}

type EndingView content.EndingCopy

func Ending(ending logic.Ending, cont *content.Content) EndingView {
	return EndingView(cont.Endings[ending])
}
