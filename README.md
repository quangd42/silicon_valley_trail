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

<!--toc:end-->

## Demo

- Screen recording / hosted demo: `TODO`

## Quick Start

Requirements:

- Go `1.26.1` installed

Instructions:

```bash
# Clone the repo
git clone git@github.com:quangd42/silicon_valley_trail.git svt && cd svt

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
- `internal/weather`: live API client + mock weather service
- `internal/save`: save/load, currently only supports JSON format
- `internal/content`: route data and game copy
- `internal/view`: view models, built from game copy and game state
- `internal/ui`: rendering using view models, currently only supports terminal CLI

Explanation:

`program` runs the main game loop:

- `program` renders game content and player interactions using methods from `ui`, such as `RenderDay` and `PromptSelection`.
- `PromptSelection` offers player options based on the current game state and collects input.
- `program` processes input, modifies game state, and evaluates winning/losing conditions with `logic`.
- `program` also pulls in external data from `weather` and provides it to `logic` as a gameplay variable.
- `program` uses `save` when the player selects `Save Game` or `Load Game`.

`model` defines all data types for the game.

`ui` is built to depend on `view` and `content` for a structured game view, so that other renderers (such as a web UI) can be plugged
in. `content` allows game copy to be updated separately from gameplay logic and display.

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
- Improving unit test coverage (edge cases), and unit test stub generation.
- Code reviewing.
- General copy editing: checking spelling and grammar, event copy generating.

## Design Notes

### Game loop & balance approach

The core resource is `Product Readiness`. When the player reaches the final destination, the game rolls a number in [0-100); if `Product Readiness` is not less than the rolled value, the player wins, otherwise they lose. The main pressure of the game comes from choosing the right series of actions to improve these odds while managing resources to avoid losing conditions, such as running out of cash or going without coffee for too long. The game is balanced around this pressure.

Each game loop is a day:

- Choose one action each day
- Manage resources `Cash`, `Morale`, `Coffee`, `Hype`, and `Product Readiness`
- The loop ends when the player has traveled across 10 real Bay Area locations

Current actions:

- Travel to the next location (costs cash, coffee, and morale)
- Rest and recover (restore morale and coffee, costs cash)
- Work on product (increase product readiness, costs coffee and morale)
- Marketing push (increase hype, costs cash and coffee)

Win / lose conditions:

- Win by reaching San Francisco and resolving the final pitch
- Lose if cash reaches zero
- Lose if coffee stays at zero for too long

### Live weather data

At each location, weather data modifies action outcomes and incentivizes or disincentivizes certain choices:

- `Clear`: boosts travel / build outcomes
- `Rainy`: makes travel harder and weakens marketing
- `Fog`: hurts travel and build output
- `Cloudy`: helps rest and slightly improves marketing

Weather data is fetched live from `WeatherAPI.com`. A mock fallback is provided for offline play and local development.
The mock fallback is also used when a live request fails.

### Data models

`State` holds game state data: `Day`, `Location`, `Resources`, `Weather`, and other game-check variables. In each game session,
`State` is modified as the game progresses. When the game is saved, data from `State` is serialized to disk as a JSON file.
When a saved game is loaded, `State` is populated from a JSON file on disk.

`Location` represents each location in the trail. It carries coordinates as input for `WeatherAPI`.
Each field in `Resources` is a simple integer representing either a percentage (`Morale`, `Hype`, `Product`) or a raw count (`Cash`, `Coffee`).
`Weather`, `Action` and `Control` are all enum-like types:

- `Weather`: weather condition on the day
- `Action`: available game choice
- `Control`: game session control (save, load, quit)

### Tradeoffs / If I had more time:

- Prefetch the next location's weather asynchronously so timeout handling (and fallback decisions) are done before the player advances to the next location, to
  ensure no blank screen.
- One time events versus recurring events: some events should be removed from the pool, while some others can reoccur based on certain conditions.
- Recruitment & trust mechanics. Recruits improve build velocity with diminishing returns. Trust and morale determine whether new team members stay. This also enables an alternative emotional ending `EndingAlone`, which occurs when the team size is 1 when the game is won (the team size is 2 at the start of the game: the player and buddy Pete).
- Multiple game saves instead of just 1. Explore using SQLite instead of JSON as game save storage format.
- Data schema versioning: make sure old game saves will still work if later versions have new state.
- General game balancing through action costs, weather effects, or new factors such as:
  - Traveling cost that varies proportionally to distance between locations, populated by Google Maps API.
  - Inventory system, with items acquired through events. Items provide modifier effects and enable trading.
