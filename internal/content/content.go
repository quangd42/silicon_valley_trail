package content

import "github.com/quangd42/silicon_valley_trail/internal/model"

type Content struct {
	Intro  string
	Route  []model.Location
	Ending string
}

func Load() *Content {
	return &Content{
		Intro:  introCopy(),
		Route:  DefaultRoute(),
		Ending: endingCopy(),
	}
}

func introCopy() string {
	return `Welcome to Silicon Valley Trail!

You and your best bud Pete set out from your HQ in San Jose to San Francisco to attend a major investor meeting. Your product: a sleeping mask that lets people relive childhood memories through dreams.

Will you be able to impress the investors?

`
}

func endingCopy() string {
	return `Congratulations! You reached San Francisco! After a hasty coffee you went into the meeting room where investors are waiting for you.

The presentation went well. You got a verbal contract on the spot...

Your team is estatic!

...

But what did it cost?...`
}
