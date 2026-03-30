package view

import "github.com/quangd42/silicon_valley_trail/internal/content"

type IntroView []string

func Intro(cont *content.Content) IntroView {
	return IntroView(cont.Intro)
}

type EndingView string

func Ending(cont *content.Content) EndingView {
	return EndingView(cont.Ending)
}
