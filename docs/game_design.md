# Silicon Valley Trail

You and your best bud Pete set out from your HQ in San Jose to San Francisco to attend a major investor meeting. Your product: a sleeping mask that lets people relive childhood memories through dreams.

Will you be able to impress the investors?

## High level concept

A startup road trip sim game, taking inspiration from Undertale: choices not only impact **whether** you'll arrive, but also **how** you'll arrive.

## 1. Game loop
1. Show current day, location, resources, party, and weather.
2. Player chooses one action.
3. Resolve action cost/benefit.
4. If the player chose `Travel`, move to the next location.
5. Trigger one arrival event.
6. Apply event choice outcomes.
7. Check lose conditions.
8. If at San Francisco, resolve ending.

## 2. Resources & state

### Visible
- `Cash`
  - used for travel, food, lodging, marketing, and event costs
- `Morale`
  - tracks team energy and mood
- `Coffee`
  - short-term stamina resource
- `Hype`
  - public excitement and investor attention
- `Product Readiness`
  - how convincing and stable the startup looks by the final pitch

### Meta
- `Trust`
  - whether the team still believes in your leadership
  - used for companion retention and ending text

## 3. Map

Fixed route:
1. San Jose
2. Santa Clara
3. Sunnyvale
4. Mountain View
5. Palo Alto
6. Menlo Park
7. Redwood City
8. San Mateo
9. South San Francisco
10. San Francisco

## 4. Events & choices

### Starting State
- start in San Jose
- party: founder + Pete
- cash: medium-low
- morale: healthy
- coffee: limited
- hype: low-medium
- product readiness: low-medium
- trust: healthy

Exact value to be tuned to balance the game.

### Core Actions
1. Travel
- move to next location
- costs cash and coffee
- weather may add extra penalties

2. Rest
- restore morale
- may restore a little coffee
- costs cash
- slightly helps trust when the team is strained

3. Work On Product
- increases product readiness
- costs coffee
- small morale cost

4. Marketing Push
- increases hype
- costs cash
- may reduce trust if overused or if the product is weak

### Events
One event happens after each movement.

Event categories:
- startup opportunity
- supply / travel problem
- weather-related complication
- recruitable teammate encounter
- team conflict moment

Each event should offer 2-3 choices with real tradeoffs:
- gain cash, lose trust
- gain hype, lose product readiness
- protect morale, lose time or money

### Party
#### Pete
- always starts with you
- emotional anchor of the game
- main signal for whether the founder is still leading like a person, not just chasing the pitch

#### Additional teammates
- max 3 recruitable teammates
- increase product readiness improvement rate
- a simple leave/stay rule driven by trust

### Endings
#### Win
Reach San Francisco. Variants are triggered by ending state.

Variants:
- `Together`
  - Pete is still with you
  - Meeting result is a win

- `The real treasure`
  - Pete is still with you
  - Meeting result is a nope

- `Alone`
  - you reach San Francisco, but Pete and the rest of the team are gone
  - Meeting result good or bad

The investor meeting result can be influenced by:
- product readiness
- hype +- 30%
- weather +- 5%

#### Lose
- cash reaches zero and you cannot pay a required cost
- morale collapses
- coffee stays at zero for 2 days
- optional: extreme trust collapse causes Pete to leave, but not instant game over

## 5. Public Web API integration (at least 1 required)

Use OpenWeatherMap with env-based API key loading and mock fallback.

Weather affects gameplay, not just flavor:
- `Rain`: travel uses extra coffee, morale down
- `Heat`: morale down unless resting
- `Fog`: travel stress, small morale penalty
- `Clear / Cloudy`: neutral

