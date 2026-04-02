package gamedef

import (
	"github.com/quangd42/silicon_valley_trail/internal/logic"
	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type EventData struct {
	Name         string
	Narrative    Narrative
	ChoicesLabel string
	Choices      []EventChoiceData
}

type EventChoiceData struct {
	Name      string
	Desc      string
	Narrative Narrative
	Effect    logic.Effect
}

func eventData() []EventData {
	return []EventData{
		{
			Name: "We meet again!",
			Narrative: Narrative{
				"A cheery disheveled fellow approaches you gleefully. You do not know this man.",
				`"It's me, Ranwid! Have any goods for me today? The usual? A fella like me can't
make it alone, you know?"`,
				"You eye him suspiciously and consider your options...",
			},
			ChoicesLabel: "Give him...",
			Choices: []EventChoiceData{
				{
					Name: "some coffee.",
					Desc: "-2 Coffee, +10 Hype",
					Narrative: Narrative{
						`Ranwid: "Exquisite! Was feeling parched."

*Glup glup glup*

He downs the coffee in one go and lets out a satisfied burp.`,
						`He wipes his mouth, climbs onto a bench, and starts shouting at passing commuters.

Ranwid: "Listen up! Dream-mask startup! Future unicorn! Tell your coworkers!"

By the time he disappears into the crowd, your mentions are on fire.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +10,
							},
						}
					},
				},
				{
					Name: "some money.",
					Desc: "-200 Cash, +10 Hype",
					Narrative: Narrative{
						`Ranwid: "Magnificent! This will be quite handy if I run into those *mask wearing hoodlums* again."`,
						`He pockets the cash, gives you a crooked salute, and vanishes into the station.

Minutes later, your phone starts buzzing. Ranwid has apparently spent the whole sum boosting posts
and bribing strangers to wear your sticker.

Absurd, but effective.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -200,
								Hype: +10,
							},
						}
					},
				},
				{
					Name: "ignore him",
					Desc: "...",
					Narrative: Narrative{
						`Ranwid: "Aaaaagghh!! What a jerk you are sometimes!"

He runs away.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{}
					},
				},
			},
		},
		{
			Name: "Brutal Honesty",
			Narrative: Narrative{
				"A founder in an immaculate camel coat glances at your demo over your shoulder.",
				`"Too many words," she says. "Too many buttons too. I can fix that. For a fee, obviously."`,
				"Pete mutters that this is either your lucky break or a highly efficient scam.",
			},
			ChoicesLabel: "Let her...",
			Choices: []EventChoiceData{
				{
					Name: "tear apart the demo.",
					Desc: "-10 Morale, +8 Product",
					Narrative: Narrative{
						`She strips jargon, cuts whole flows, and forces you to explain the product like a normal
human being.

It is humiliating.`,
						`It is also the best your demo has ever looked, which somehow makes it worse.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale:  -10,
								Product: +8,
							},
						}
					},
				},
				{
					Name: "post about you.",
					Desc: "-150 Cash, +8 Hype",
					Narrative: Narrative{
						`You wire over the consulting fee.

An hour later, she has posted a glowing thread about your company, and your notifications finally
start moving.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -150,
								Hype: +8,
							},
						}
					},
				},
				{
					Name: "find someone else to insult.",
					Desc: "-2 Morale",
					Narrative: Narrative{
						`You thank her for the offer and keep moving.

Pete says the product should probably survive one less opinion today.

You still spend the next block wondering if she was right.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Warehouse Demo Night",
			Narrative: Narrative{
				"Down an alley lit by cheap LEDs, a demo night is spinning up in a converted warehouse.",
				`"One startup bailed," the organizer says. "You want the slot or not?"`,
				"Extension cords snake across the floor. So do opportunities.",
			},
			ChoicesLabel: "Take the slot and...",
			Choices: []EventChoiceData{
				{
					Name: "demo the prototype live.",
					Desc: "-6 Coffee, +5 Product, +4 Hype",
					Narrative: Narrative{
						`You scramble through setup, patch one last bug, and somehow make it onstage in time.

The room actually leans in.`,
						`The questions afterward are even better than the applause.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -6,
								Product: +5,
								Hype:    +4,
							},
						}
					},
				},
				{
					Name: "turn it into a stunt.",
					Desc: "-400 Cash, +12 Hype",
					Narrative: Narrative{
						`You spend fast on lights, pizza, and a ridiculous vinyl banner.

By midnight, half the warehouse knows your company name, and the other half is filming it.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -400,
								Hype: +12,
							},
						}
					},
				},
				{
					Name: "pass.",
					Desc: "-2 Hype",
					Narrative: Narrative{
						`You decide chaos can be somebody else's growth strategy tonight.

Unfortunately, somebody else takes the slot and gets the buzz instead.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Hype: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Basement Hackspace",
			Narrative: Narrative{
				"A basement hackspace offers one open desk, one rattling fan, and a pot of coffee old enough to vote.",
				`The manager shrugs. "Use the night well."`,
			},
			ChoicesLabel: "Spend the night to...",
			Choices: []EventChoiceData{
				{
					Name: "ship one more feature.",
					Desc: "-3 Coffee, +6 Product",
					Narrative: Narrative{
						`You code under fluorescent buzz until your eyes go grainy.

At sunrise, the build is cleaner, faster, and far less embarrassing.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -3,
								Product: +6,
							},
						}
					},
				},
				{
					Name: "host a tiny launch party.",
					Desc: "-150 Cash, +8 Hype",
					Narrative: Narrative{
						`You buy snacks, slap stickers on everything, and stream a shameless little midnight launch.

It is messy. It is loud. It works.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -150,
								Hype: +8,
							},
						}
					},
				},
				{
					Name: "rent the nap pod.",
					Desc: "-80 Cash, +12 Morale, +3 Coffee",
					Narrative: Narrative{
						`You pay for the pod, collapse instantly, and wake up to a paper cup of terrible coffee on
the floor outside.

You have never felt richer.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:   -80,
								Morale: +12,
								Coffee: +3,
							},
						}
					},
				},
			},
		},
		{
			Name: "Canceled Lunch",
			Narrative: Narrative{
				"Outside a hotel, a flustered assistant mistakes you for another founder and hands you a canceled lunch reservation with a junior VC.",
				`"They're already seated. If you want the table, walk in like you belong there."`,
			},
			ChoicesLabel: "Use the table to...",
			Choices: []EventChoiceData{
				{
					Name: "sell the vision.",
					Desc: "-2 Coffee, +8 Hype",
					Narrative: Narrative{
						`You walk in like you belong there and pitch with the kind of intensity that empties a cup
before the appetizers arrive.

The junior VC does not commit to anything.`,
						`They do promise to remember your name.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +8,
							},
						}
					},
				},
				{
					Name: "ask for hard feedback.",
					Desc: "-6 Morale, +6 Product",
					Narrative: Narrative{
						`You skip the performance and ask what is broken.

They tell you.`,
						`Pete goes quiet, you take notes, and the product comes out sharper for it.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale:  -6,
								Product: +6,
							},
						}
					},
				},
				{
					Name: "sell the reservation.",
					Desc: "+120 Cash, -6 Morale",
					Narrative: Narrative{
						`Another founder offers cash on the spot for the table.

You take the deal, then spend the walk convincing yourself this was "capital efficient."`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:   +120,
								Morale: -6,
							},
						}
					},
				},
			},
		},
		{
			Name: "Open Lab",
			Narrative: Narrative{
				"A community college lab is still open, and a handful of students are poking at half-finished side projects.",
				"One of them notices your prototype and asks if they can try it.",
			},
			ChoicesLabel: "Use the lab to...",
			Choices: []EventChoiceData{
				{
					Name: "stress-test the prototype.",
					Desc: "-2 Coffee, +5 Product",
					Narrative: Narrative{
						`Within minutes, the students have broken your onboarding in three different ways.

By the time they leave, your bug list is longer and your prototype is much better.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -2,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "shoot a social clip.",
					Desc: "-5 Morale, +8 Hype",
					Narrative: Narrative{
						`You redo the demo again and again until the students finally get the shot they want.

It is exhausting.`,
						`It is also the first time your product has looked cool on camera.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -5,
								Hype:   +8,
							},
						}
					},
				},
				{
					Name: "share a late meal.",
					Desc: "-120 Cash, +10 Morale",
					Narrative: Narrative{
						`You buy a round of late-night noodles and end up talking shop for an hour.

By the time you leave, your team feels human again.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:   -120,
								Morale: +10,
							},
						}
					},
				},
			},
		},
		{
			Name: "Podcast Booth",
			Narrative: Narrative{
				"A tiny startup podcast booth is recording outside a coffee shop, and the host spots your badge.",
				`"You building something weird enough for radio?" she asks.`,
				"Pete says that is either an insult or an invitation.",
			},
			ChoicesLabel: "Step up and...",
			Choices: []EventChoiceData{
				{
					Name: "do the interview.",
					Desc: "-4 Morale, +8 Hype",
					Narrative: Narrative{
						`You smile on command, compress the company into three confident sentences, and somehow make it sound intentional.`,
						`By the end, you are exhausted, but memorable.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -4,
								Hype:   +8,
							},
						}
					},
				},
				{
					Name: "ask what confused her.",
					Desc: "-2 Coffee, +5 Product",
					Narrative: Narrative{
						`She points out exactly where your explanation turns to mush, then makes you say it again in plain English.`,
						`It stings, but the pitch gets sharper fast.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -2,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "pretend you have another meeting.",
					Desc: "-2 Morale",
					Narrative: Narrative{
						`You nod politely, check a calendar notification that does not exist, and keep walking.

The move works, but only if you do not count your own opinion.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Coworking Day Pass",
			Narrative: Narrative{
				"A coworking space is offering one-day passes to anyone who looks sleep-deprived enough to need them.",
				`The receptionist gestures at the lounge. "Pick a corner and try not to pivot loudly."`,
			},
			ChoicesLabel: "Use it to...",
			Choices: []EventChoiceData{
				{
					Name: "book the tiny meeting room.",
					Desc: "-140 Cash, +7 Product",
					Narrative: Narrative{
						`You take the whiteboard, close the glass door, and finally work through the ugliest part of the demo flow.`,
						`Privacy helps more than you expected.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:    -140,
								Product: +7,
							},
						}
					},
				},
				{
					Name: "raid the kitchen and reset.",
					Desc: "-70 Cash, +3 Coffee, +6 Morale",
					Narrative: Narrative{
						`You overpay for the pass, then absolutely get your money's worth in stale cold brew and suspicious granola.`,
						`Pete looks human again by the time you leave.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:   -70,
								Coffee: +3,
								Morale: +6,
							},
						}
					},
				},
				{
					Name: "just charge the laptop.",
					Desc: "...",
					Narrative: Narrative{
						`You borrow an outlet, ignore the networking energy, and move on once the battery is no longer flashing red.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{}
					},
				},
			},
		},
		{
			Name: "Founder Support Thread",
			Narrative: Narrative{
				"An alumni founder group chat is melting down over onboarding metrics, API outages, and whether anyone has seen the office HDMI adapter.",
				"Your phone will not stop vibrating.",
			},
			ChoicesLabel: "Jump in and...",
			Choices: []EventChoiceData{
				{
					Name: "answer the thread.",
					Desc: "-2 Coffee, +7 Hype",
					Narrative: Narrative{
						`You drop advice, screenshots, and one alarmingly honest anecdote about your own launch mistakes.`,
						`By the end of the hour, more people know your company name than before.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +7,
							},
						}
					},
				},
				{
					Name: "ask them to break the prototype.",
					Desc: "-4 Morale, +5 Product",
					Narrative: Narrative{
						`They comply with unsettling enthusiasm and return a list of bugs long enough to hurt your feelings.`,
						`Still, the product gets sturdier because of it.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale:  -4,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "mute the thread.",
					Desc: "-2 Hype",
					Narrative: Narrative{
						`You press mute, feel briefly powerful, and put the phone face down.

By the time you check back in, three other founders have filled the silence.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Hype: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Sponsor Leftovers",
			Narrative: Narrative{
				"A tech mixer is tearing down nearby, and several sponsor tables have been abandoned in a state between generosity and neglect.",
				"Nobody seems eager to inventory the leftovers.",
			},
			ChoicesLabel: "Make use of it and...",
			Choices: []EventChoiceData{
				{
					Name: "grab the energy drinks.",
					Desc: "+3 Coffee, -3 Morale",
					Narrative: Narrative{
						`The cans are warm, aggressively branded, and almost certainly bad for you.`,
						`They still get the job done.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: +3,
								Morale: -3,
							},
						}
					},
				},
				{
					Name: "claim the abandoned banner stand.",
					Desc: "-80 Cash, +6 Hype",
					Narrative: Narrative{
						`You spend a little to print something presentable, then mount it on the nicest stand anyone forgot to carry out.`,
						`It makes your whole operation look more legitimate than it is.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -80,
								Hype: +6,
							},
						}
					},
				},
				{
					Name: "sort through the cable pile.",
					Desc: "-2 Morale, +4 Product",
					Narrative: Narrative{
						`You crawl through adapters, dongles, and unlabeled chargers until you finally find the one thing your setup was missing.`,
						`It is tedious, but useful.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale:  -2,
								Product: +4,
							},
						}
					},
				},
			},
		},
		{
			Name: "Microgrant Table",
			Narrative: Narrative{
				"A nonprofit booth is handing out flyers for tiny founder grants with huge application forms attached.",
				`"If you can survive the paperwork," the volunteer says, "you can survive anything."`,
			},
			ChoicesLabel: "Sit down and...",
			Choices: []EventChoiceData{
				{
					Name: "fill out the application.",
					Desc: "-7 Morale, +150 Cash",
					Narrative: Narrative{
						`You spend an annoying amount of time describing your "impact story" to a form field that clearly hates you.`,
						`Against the odds, the grant comes through.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -7,
								Cash:   +150,
							},
						}
					},
				},
				{
					Name: "hand over the pitch deck.",
					Desc: "-2 Coffee, +6 Hype",
					Narrative: Narrative{
						`The volunteer cannot fund you, but they can absolutely post you into the orbit of several local founder circles.`,
						`That turns out to be worth something.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +6,
							},
						}
					},
				},
				{
					Name: "avoid the line.",
					Desc: "-2 Morale",
					Narrative: Narrative{
						`You decide there are only so many forms a person can survive in one week.

The volunteer watches you leave with visible disappointment.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Parking Lot Pilot",
			Narrative: Narrative{
				"A rideshare queue has stalled long enough for several commuters to become curious about anything that breaks the boredom.",
				"One of them points at your prototype and asks if it is real.",
			},
			ChoicesLabel: "Use the moment and...",
			Choices: []EventChoiceData{
				{
					Name: "let commuters try the demo.",
					Desc: "-2 Coffee, +7 Hype",
					Narrative: Narrative{
						`You run fast, improvised demos between car arrivals and somehow turn the whole thing into a tiny spectacle.`,
						`At least a few people leave talking about you.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +7,
							},
						}
					},
				},
				{
					Name: "watch where they get stuck.",
					Desc: "-4 Morale, +5 Product",
					Narrative: Narrative{
						`The testers immediately find the exact parts of the flow you were secretly hoping nobody would notice.`,
						`Embarrassing, but extremely useful.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale:  -4,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "pack up early.",
					Desc: "-2 Hype",
					Narrative: Narrative{
						`You decide not every accidental audience needs to become a growth experiment.

Someone else seizes the moment before the queue starts moving again.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Hype: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Office Liquidation",
			Narrative: Narrative{
				"A startup two blocks over is shutting down and liquidating half an office worth of gear at folding-table prices.",
				"The atmosphere is grim, but the hardware is real.",
			},
			ChoicesLabel: "Walk out with...",
			Choices: []EventChoiceData{
				{
					Name: "the spare monitor.",
					Desc: "-150 Cash, +7 Product",
					Narrative: Narrative{
						`The extra screen makes bug fixing and demo prep noticeably less painful the moment you set it up.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:    -150,
								Product: +7,
							},
						}
					},
				},
				{
					Name: "a little moving money.",
					Desc: "-4 Morale, +100 Cash",
					Narrative: Narrative{
						`You help carry dead desktops and sad plants down three flights of stairs for cash and a handshake neither side enjoys.`,
						`Useful money, bleak vibes.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -4,
								Cash:   +100,
							},
						}
					},
				},
				{
					Name: "nothing but perspective.",
					Desc: "...",
					Narrative: Narrative{
						`You take the hint, grip your backpack a little tighter, and head back outside.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{}
					},
				},
			},
		},
		{
			Name: "Quiet Hotel Lobby",
			Narrative: Narrative{
				"A hotel lobby near the conference district is so quiet that even your bad ideas sound expensive.",
				"The chairs are deep, the wifi is fast, and nobody has asked if you belong here yet.",
			},
			ChoicesLabel: "Use the calm to...",
			Choices: []EventChoiceData{
				{
					Name: "revise the deck by the outlet.",
					Desc: "-2 Coffee, +5 Product",
					Narrative: Narrative{
						`You camp by an outlet, fix three ugly slides, and make the whole story feel cleaner.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -2,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "buy overpriced tea and decompress.",
					Desc: "-60 Cash, +8 Morale",
					Narrative: Narrative{
						`The tea costs too much and tastes faintly of lemons and regret, but the break helps.`,
						`You leave steadier than you arrived.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:   -60,
								Morale: +8,
							},
						}
					},
				},
				{
					Name: "borrow the wifi and move on.",
					Desc: "...",
					Narrative: Narrative{
						`You answer a few emails fast, avoid eye contact with the concierge, and keep walking.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{}
					},
				},
			},
		},
		{
			Name: "Community QA Night",
			Narrative: Narrative{
				"A local builder meetup is running an informal QA night, and someone waves you toward an empty whiteboard.",
				`"If it breaks," they promise, "we'll tell you exactly how."`,
			},
			ChoicesLabel: "Join in and...",
			Choices: []EventChoiceData{
				{
					Name: "run bug triage.",
					Desc: "-2 Coffee, +5 Product",
					Narrative: Narrative{
						`You leave with a longer checklist, a better sense of what matters, and a much more reliable flow.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -2,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "raffle off stickers and collect emails.",
					Desc: "-90 Cash, +6 Hype",
					Narrative: Narrative{
						`You spend a little on cheap swag, then turn the whole thing into a tiny mailing-list heist.`,
						`It works better than it should.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -90,
								Hype: +6,
							},
						}
					},
				},
				{
					Name: "thank everyone and leave.",
					Desc: "-2 Morale",
					Narrative: Narrative{
						`You decide not every useful room needs to become your room tonight.

It is sensible, but it still feels like chickening out.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -2,
							},
						}
					},
				},
			},
		},
		{
			Name: "Train Delay",
			Narrative: Narrative{
				"A departure board flickers red and buys everyone on the platform twenty irritated extra minutes.",
				"Dead time is still time if you are willing to be annoying about it.",
			},
			ChoicesLabel: "Use the delay to...",
			Choices: []EventChoiceData{
				{
					Name: "rewrite the opening slide.",
					Desc: "-2 Coffee, +4 Product",
					Narrative: Narrative{
						`You rework the first impression of the pitch while standing under a broken speaker and somehow make it better.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -2,
								Product: +4,
							},
						}
					},
				},
				{
					Name: "pitch stranded commuters.",
					Desc: "-3 Morale, +6 Hype",
					Narrative: Narrative{
						`Most people ignore you, a few laugh, and one person posts the whole thing with a surprisingly kind caption.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Morale: -3,
								Hype:   +6,
							},
						}
					},
				},
				{
					Name: "buy snacks and wait it out.",
					Desc: "-50 Cash, +2 Coffee, +6 Morale",
					Narrative: Narrative{
						`You spend a little, split the snacks with Pete, and let yourselves be people for ten whole minutes.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:   -50,
								Coffee: +2,
								Morale: +6,
							},
						}
					},
				},
			},
		},
		{
			Name: "Noodle Shop Advice",
			Narrative: Narrative{
				"A noodle shop owner recognizes the startup panic in your eyes before you even sit down.",
				`"You want food, or you want advice?" she asks.`,
			},
			ChoicesLabel: "Order and...",
			Choices: []EventChoiceData{
				{
					Name: "trade a meal for blunt feedback.",
					Desc: "-100 Cash, +5 Product",
					Narrative: Narrative{
						`She listens to the pitch, slurps once, and points out the part that still sounds fake.`,
						`The food helps. The honesty helps more.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash:    -100,
								Product: +5,
							},
						}
					},
				},
				{
					Name: "let her talk you up to the room.",
					Desc: "-80 Cash, +6 Hype",
					Narrative: Narrative{
						`Within minutes, half the shop knows your company name and the other half thinks they discovered it first.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Cash: -80,
								Hype: +6,
							},
						}
					},
				},
				{
					Name: "eat in peace.",
					Desc: "...",
					Narrative: Narrative{
						`You choose silence, broth, and ten quiet minutes without optimizing anything.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{}
					},
				},
			},
		},
		{
			Name: "Photo Wall",
			Narrative: Narrative{
				"A retail pop-up nearby has a temporary photo wall, decent lighting, and absolutely no opinion about what gets filmed there.",
				"Pete is already checking camera angles.",
			},
			ChoicesLabel: "Use the setup to...",
			Choices: []EventChoiceData{
				{
					Name: "film a polished demo clip.",
					Desc: "-2 Coffee, +6 Hype",
					Narrative: Narrative{
						`It takes more takes than your dignity would prefer, but the final clip makes the product look real in a way screenshots never do.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee: -2,
								Hype:   +6,
							},
						}
					},
				},
				{
					Name: "record a cleaner onboarding flow.",
					Desc: "-2 Coffee, +4 Product",
					Narrative: Narrative{
						`You replay the first-run experience until the clumsy parts are impossible to ignore, then fix them.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Coffee:  -2,
								Product: +4,
							},
						}
					},
				},
				{
					Name: "leave it for someone else.",
					Desc: "-2 Hype",
					Narrative: Narrative{
						`You decide the internet can survive one less video of your face today.

Unfortunately, the internet rewards the founder who did not hesitate.`,
					},
					Effect: func(_ *model.State, _ logic.Context) logic.Change {
						return logic.Change{
							Delta: model.Resources{
								Hype: -2,
							},
						}
					},
				},
			},
		},
	}
}
