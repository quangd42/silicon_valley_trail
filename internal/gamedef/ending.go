package gamedef

import "github.com/quangd42/silicon_valley_trail/internal/logic"

type EndingCopy struct {
	Narrative Narrative
	Explain   string
}

func endingCopy() map[logic.Ending]EndingCopy {
	return map[logic.Ending]EndingCopy{
		logic.EndingNone: {}, // This is just a placeholder, this is not a real ending
		logic.EndingNoCash: {
			Narrative: Narrative{
				"You count what is left in the company account, look at Pete, and both of you reach the same conclusion at once.",
				"San Francisco is not happening today.",
				"You have just enough left to get home and start figuring out what comes next.",
			},
			Explain: "You ran out of cash before reaching the meeting.",
		},
		logic.EndingNoCoffee: {
			Narrative: Narrative{
				"The first coffee-less day is survivable. The second turns the roadmap into abstract art.",
				"By the time either of you can form a sentence again, the run is over.",
			},
			Explain: "You went too long without coffee.",
		},
		logic.EndingNoProductFit: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco.",
				"You step into the investor meeting with a story, a prototype, and just enough confidence to survive the opening slides.",
				"The questions come fast. The weak spots come faster.",
				"Maybe getting here was the easy part.",
			},
			Explain: "You reached the destination, but the product and buzz still need work.",
		},
		logic.EndingMomentum: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco.",
				"The meeting is tense, but the demo holds and the room actually leans in.",
				"You do not leave with certainty, but you do leave with follow-up emails, sharper questions, and a little real momentum.",
			},
			Explain: "You arrived with enough product and hype to keep the company moving.",
		},
		logic.EndingPerfection: {
			Narrative: Narrative{
				"Congratulations! You made it to San Francisco.",
				"For one brief, unnatural afternoon, everything works: the story is clear, the demo behaves, and even the skeptical questions turn into interest.",
				"Nobody calls it perfect out loud, but it is close enough to feel dangerous.",
			},
			Explain: "You arrived with a breakout product and undeniable hype.",
		},
		logic.EndingAlone: {}, // To be added if we ever get there
	}
}
