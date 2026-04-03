# Silicon Valley Trail

<!--toc:start-->

- [Demo](#demo)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
  - [Configure live weather](#configure-live-weather)
- [Architecture](#architecture)
  - [Dependencies](#dependencies)
- [Tests](#tests)
- [AI Usage](#ai-usage)
- [Design Notes](#design-notes)
  - [Game loop & balance approach](#game-loop-balance-approach)
  - [Live weather data](#live-weather-data)
  - [Data models](#data-models)
  - [Tradeoffs / If I had more time:](#tradeoffs-if-i-had-more-time)
- [Example Interaction](#example-interaction)

<!--toc:end-->

## Demo

(YouTube link, opens in new tab)
<a href="https://www.youtube.com/watch?v=HWU_JS50m4I" target="_blank" rel="noopener noreferrer">
<img src="https://img.youtube.com/vi/HWU_JS50m4I/maxresdefault.jpg" alt="Demo">
</a>

## Quick Start

Requirements:

- Go `1.26.1` installed

Instructions:

```bash
# Clone the repo
git clone https://github.com/quangd42/silicon_valley_trail.git svt && cd svt

# Run with mock weather (default if no API key is set)
go build .
./silicon_valley_trail
# Alternatively
go run .

# To run with live weather data, see Configuration below
```

## Configuration

The game reads configuration from environment variables. It also loads variables from a local `.env` file.

Example ([.env.example](.env.example)):

```env
WEATHERAPI_KEY=your_key_here
WEATHERAPI_MOCK=true
WEATHERAPI_TIMEOUT_MS=3000
SAVE_PATH=custom_svt_save.json
```

Notes:

- If `WEATHERAPI_KEY` is empty, the game defaults to mock weather.
- `WEATHERAPI_MOCK=true` forces mock mode even if a key is present.
- `SAVE_PATH` is optional; default is `svt_save.json`.

### Configure live weather

To get a `WeatherAPI` key, sign up for an account on <https://www.weatherapi.com/>. The API key is available on the dashboard
at <https://www.weatherapi.com/my/> after signing in.

Then provide your API key in the `.env` file:

```env
WEATHERAPI_KEY=your_key_here
WEATHERAPI_MOCK=false
```

Alternatively, run the game directly with the key in the environment:

```bash
WEATHERAPI_KEY=your_key_here go run .
```

## Architecture

High-level package layout:

- `main.go`: main entry point, wires game dependencies
- `internal/program`: game loop and flow control
- `internal/logic`: action resolution, resource deltas, endings
- `internal/model`: core game data types
- `internal/gamedef`: authored route, action, weather, event, and ending data
- `internal/weather`: live API client + mock weather service
- `internal/save`: save/load, currently only supports JSON format
- `internal/view`: view models, built from authored game data and game state
- `internal/ui`: rendering using view models, currently only supports terminal CLI

Explanation:

`program` runs the main game loop:

- `program` refreshes weather for the current location, renders the day state, collects player input, and resolves the chosen action.
- After each travel action, `program` triggers exactly one arrival event before the next day continues.
- Event selection prefers weather-specific event pools first, then falls back to the general pool, and removes events from the pool once used.
- `program` evaluates losing conditions with `logic`, and resolves the final meeting ending tier once the route is complete.
- `program` uses `save` to persist and restore the full game state, including any in-progress event.

`model` defines the persisted game state, including resources, route progress, remaining event pools, and the current in-progress event if a save happens mid-event.

`gamedef` holds authored game content and mechanics configuration: actions, weather effects, route locations, event definitions, and ending copy.

`ui` depends on `view` so that rendering stays separate from gameplay rules and state transitions.

`save` offers the interface to persist game state to disk. It also provides one concrete implementation `JSONSaver`, which serializes
game state to a JSON file.

### Dependencies

Runtime:

- Go `1.26.1`
- `github.com/joho/godotenv`
- `golang.org/x/text`

External service:

- `WeatherAPI.com` for live weather data

## Tests

To run tests:

```bash
go test ./...
```

## AI Usage

AI was used in the creation of the code in the following ways:

- A research tool in place of Google: researching weather API providers, finding real Silicon Valley locations for the trail, and checking available features in the Go standard library.
- General project planning and brainstorming.
- Improving unit test coverage (edge cases).
- Code reviewing.
- General copy editing: generating event copy, checking spelling and grammar.

## Design Notes

### Game loop & balance approach

The core tension is survival versus polish. Reaching San Francisco ends the run. `Product Readiness` and `Hype` determine the flavor of the final investor meeting through the score `Product Readiness + Hype/2`.

Ending thresholds:

- `< 50`: `No Product Fit`
- `50-89`: `Momentum`
- `>= 90`: `Perfection`

Cash depletion and going too long without coffee are the only hard failure states. The game is balanced so that travel and optional progress actions both consume meaningful resources, and arrival events sometimes let the player trade morale or hype back into cash or coffee when survival gets tight.

In an earlier iteration, the final meeting used the score as a probability threshold to determine whether the player received a favorable ending. That created a stronger gameplay incentive to improve the score, but I removed it in the current version to align more literally with the assignment wording that reaching the destination should be the success condition.

Each game loop is a day:

- Choose one action each day
- Manage resources `Cash`, `Morale`, `Coffee`, `Hype`, and `Product Readiness`
- If the player traveled, resolve one arrival event before the next day starts
- The loop ends when the player reaches San Francisco at the end of a 10-location Bay Area route

Current actions:

- Travel to the next location (costs cash, coffee, and morale)
- Rest and recover (restore morale and coffee, costs cash)
- Work on product (increase product readiness, costs coffee and morale)
- Marketing push (increase hype, costs a lot of cash and some coffee)

End conditions:

- Reach San Francisco to complete the run
- Final meeting flavor is based on `Product Readiness + Hype/2`
- Lose if cash reaches zero
- Lose if coffee stays at zero for too long

### Live weather data

At each location, weather data modifies action outcomes and also affects event selection:

- `Clear`: boosts travel / build outcomes
- `Rainy`: makes travel harder and weakens marketing
- `Fog`: hurts travel and build output
- `Cloudy`: helps rest and slightly improves marketing

After each travel action, the game resolves one arrival event. The event system first checks for any weather-conditioned events for the current weather and draws from that pool before falling back to the general event pool. This makes the API influence both moment-to-moment action math and which situations the player encounters. This logic is demonstrated when the weather is `Rainy`.

Weather data is fetched live from `WeatherAPI.com`. A mock fallback is provided for offline play and local development. The mock fallback is also used when a live request fails.

The live weather client also caches weather per location for the duration of a game session, including across multiple days spent in the same city. That is an intentionally design choice: within the scope of a short run, the weather at one location is unlikely to change enough to justify extra API calls, and caching keeps gameplay more stable. This could be changed to refresh on a per-day basis or use a time-based cache.

### Data models

`State` holds game state data: `Day`, `Route`, `CurrentLocation`, `Resources`, `Party`, `Weather`, `NoCoffeeDayCount`, `EventPools`, and `CurrentEvent`. In each game session, `State` is modified as the game progresses. When the game is saved, data from `State` is serialized to disk as a JSON file. When a saved game is loaded, `State` is populated from a JSON file on disk.

`Party` currently starts as a fixed two-person team (`You` and Pete) and is persisted as part of the save state. It does not yet have gameplay systems that change party size mid-run, but it leaves room for future mechanics such as recruitment and trust.

`EventPools` stores the remaining general and weather-specific event IDs so events can be consumed without repeating indefinitely. `CurrentEvent` allows the game to resume safely if the player saves in the middle of an event prompt.

`Location` represents each location in the trail. It carries coordinates as input for `WeatherAPI`.
Each field in `Resources` is a simple integer representing either a percentage (`Morale`, `Hype`, `Product`) or a raw count (`Cash`, `Coffee`).
`gamedef.Definition.Events` is keyed by stable event IDs, while `State` stores only those IDs in pools and save data.
`Weather`, `Action` and `Control` are all enum-like types:

- `Weather`: weather condition on the day
- `Action`: available game choice
- `Control`: game session control (save, load, quit)

### Tradeoffs / If I had more time:

- Prefetch the next location's weather asynchronously so timeout handling (and fallback decisions) are done before the player advances to the next location, to ensure no blank screen.
- One time events versus recurring events: some events should be removed from the pool, while some others can reoccur based on certain conditions.
- Recruitment & trust mechanics. Recruits could improve build velocity with diminishing returns. Trust and morale could determine whether new team members stay. This could also enable an alternative emotional ending `EndingAlone`, which would occur if the team size is 1 when the game is won (the team size currently starts at 2: the player and buddy Pete).
- Multiple game saves instead of just 1. Explore using SQLite instead of JSON as game save storage format.
- Data schema versioning: make sure old game saves will still work if later versions have new state.
- Weather API issues can be logged. Currently it falls back to mock provider silently.
- General game balancing through action costs, weather effects, or new factors such as:
  - Traveling cost that varies proportionally to distance between locations, populated by Google Maps API.
  - Add more RNGs to the event outcomes.
  - Inventory system, with items acquired through events. Items provide modifier effects and enable trading.

## Example Interaction

Sample run with mock weather, trimmed for length. Weather and event draws vary between runs.
The chosen choice of each day this sample run is marked with `**choice**` for legibility.

```text
$ go run .
================================================================================
SILICON VALLEY TRAIL - Main Menu
================================================================================
1. **New Game**
2. Load Game
3. Quit Game

Enter choice (1-3): 1

Welcome to Silicon Valley Trail!

You and your best bud Pete set out from your HQ in San Jose to San Francisco for
a high-stakes investor meeting. Your product: a sleeping mask that lets people
relive childhood memories through dreams.

Reach San Francisco before the company runs out of runway.

Press Enter to begin your journey!

================================================================================
Day 1 | San Jose
Garage HQ and the start of the run.
================================================================================
Cash: $6,000 | Morale: 100% | Coffee: 26
Hype: 10% | Product Readiness: 20%
Progress: 10% to San Francisco
================================================================================
Weather: Clear
You feel productive and ready to go.
(1.Travel+  3.Build+)
--------------------------------------------------------------------------------
What will you do?
--------------------------------------------------------------------------------
1. **Travel to the next location (costs cash, coffee, and morale)**
2. Rest and recover (restore morale and coffee, costs cash)
3. Work on product (increase product readiness, costs coffee and morale)
4. Marketing push (increase hype, costs a lot of cash and some coffee)
5. Save Game
6. Quit to Menu

Enter choice (1-6): 1
...
Arrived at Santa Clara!

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
EVENT: Rain-soaked pitch
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
Rain batters the station roof in a steady metallic rhythm.
...
You decide to...
1. pitch to the stranded crowd. (-2 Morale, +8 Hype)
2. **wait it out and keep spirits up. (+4 Morale)**
3. Save Game
4. Quit to Menu

Enter choice (1-4): 2
...

================================================================================
Day 2 | Santa Clara
Corporate campuses and investor drive-bys.
================================================================================
Cash: $5,650 | Morale: 100% | Coffee: 22
Hype: 10% | Product Readiness: 20%
Progress: 20% to San Francisco
================================================================================
Weather: Rainy
It's miserable out there.
(1.Travel-  4.Marketing-)
--------------------------------------------------------------------------------
What will you do?
--------------------------------------------------------------------------------
1. Travel to the next location (costs cash, coffee, and morale)
2. Rest and recover (restore morale and coffee, costs cash)
3. **Work on product (increase product readiness, costs coffee and morale)**
4. Marketing push (increase hype, costs a lot of cash and some coffee)
5. Save Game
6. Quit to Menu

Enter choice (1-6): 3
--------------------------------------------------------------------------------
You take on the next item on the roadmap...

You're happy with the result, but everyone is tired...
--------------------------------------------------------------------------------
(Coffee -6. Morale -18%. Product +8%. Press Enter to continue...)

================================================================================
Day 3 | Santa Clara
Corporate campuses and investor drive-bys.
================================================================================
Cash: $5,650 | Morale: 82% | Coffee: 16
Hype: 10% | Product Readiness: 28%
Progress: 20% to San Francisco
================================================================================
Weather: Fog
You feel unsure about the future.
(1.Travel-  3.Build-)
--------------------------------------------------------------------------------
What will you do?
--------------------------------------------------------------------------------
1. Travel to the next location (costs cash, coffee, and morale)
2. Rest and recover (restore morale and coffee, costs cash)
3. Work on product (increase product readiness, costs coffee and morale)
4. Marketing push (increase hype, costs a lot of cash and some coffee)
5. **Save Game**
6. Quit to Menu

Enter choice (1-6): 5
--------------------------------------------------------------------------------
Game saved.
--------------------------------------------------------------------------------

================================================================================
Day 3 | Santa Clara
Corporate campuses and investor drive-bys.
================================================================================
Cash: $5,650 | Morale: 82% | Coffee: 16
Hype: 10% | Product Readiness: 28%
Progress: 20% to San Francisco
================================================================================
Weather: Fog
You feel unsure about the future.
(1.Travel-  3.Build-)
--------------------------------------------------------------------------------
What will you do?
--------------------------------------------------------------------------------
1. Travel to the next location (costs cash, coffee, and morale)
2. Rest and recover (restore morale and coffee, costs cash)
3. Work on product (increase product readiness, costs coffee and morale)
4. Marketing push (increase hype, costs a lot of cash and some coffee)
5. Save Game
6. **Quit to Menu**

Enter choice (1-6): 6

================================================================================
SILICON VALLEY TRAIL - Main Menu
================================================================================
1. New Game
2. **Load Game**
3. Quit Game

Enter choice (1-3): 2

================================================================================
Day 3 | Santa Clara
Corporate campuses and investor drive-bys.
================================================================================
Cash: $5,650 | Morale: 82% | Coffee: 16
Hype: 10% | Product Readiness: 28%
Progress: 20% to San Francisco
================================================================================
Weather: Fog
You feel unsure about the future.
(1.Travel-  3.Build-)
--------------------------------------------------------------------------------
What will you do?
--------------------------------------------------------------------------------
1. Travel to the next location (costs cash, coffee, and morale)
2. Rest and recover (restore morale and coffee, costs cash)
3. Work on product (increase product readiness, costs coffee and morale)
4. Marketing push (increase hype, costs a lot of cash and some coffee)
5. Save Game
6. Quit to Menu
```
