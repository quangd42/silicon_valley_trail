package gamedef

func introCopy() Narrative {
	return Narrative{
		`Welcome to Silicon Valley Trail!

		You and your best bud Pete set out from your HQ in San Jose to San Francisco for a high-stakes investor meeting. Your product: a sleeping mask that lets people relive childhood memories through dreams.

		Reach San Francisco before the company runs out of runway.`,
		"Instructions:\n\nArrive in San Francisco without going broke or burning out. Product readiness and hype no longer decide whether you make it there, but they will decide how impressed the investors are when you finally walk into the room.",
		`Manage your resources wisely:
		* Cash    ($): Don’t run out. No cash = the run ends early.
		* Morale  (%): Tired teams build worse and make worse calls.
		* Coffee  (ct): Your startup fuel. 2 days without it = the run ends.
		* Product (%): How convincing the product looks in the final meeting.
		* Hype    (%): Public attention. Every 2 Hype = 1 Product at the final meeting.`,
	}
}
