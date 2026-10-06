# Pokedex

A command-line Pokedex written in Go. Explore the Pokémon world through an interactive REPL: browse location areas, see which Pokémon live there, try to catch them and inspect the ones you caught.

All data comes from [PokeAPI](https://pokeapi.co). Responses are cached in memory, so revisiting a page or a Pokémon doesn't trigger a new request.

Built as part of the [Boot.dev](https://www.boot.dev) Go course.

## Requirements

- Go 1.27 or newer
- Internet access (to reach PokeAPI)

## Getting started

```sh
git clone <repository-url>
cd pokedex

# run directly
go run .

# or build a binary
go build -o pokedex
./pokedex
```

## Commands

| Command | Description |
|---|---|
| `help` | Show the available commands |
| `map` | Show the next 20 location areas |
| `mapb` | Show the previous 20 location areas |
| `explore <area>` | List the Pokémon found in a location area |
| `catch <pokemon>` | Throw a Pokéball and try to catch a Pokémon |
| `inspect <pokemon>` | Show details of a Pokémon you caught |
| `pokedex` | List every Pokémon you caught |
| `exit` | Close the Pokedex |

Input is case-insensitive.

## Example session

```text
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...
Pokedex > explore canalave-city-area
tentacool
tentacruel
staryu
...
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  -hp: 40
  -attack: 45
  ...
Types:
  -normal
  -flying
Pokedex > pokedex
Your Pokedex:
  -pidgey
```

## How catching works

Each Pokémon has a `capture_rate` from 3 to 255 (higher means easier), taken from the `/pokemon-species` endpoint. A catch rolls a 256-sided die (0–255) and succeeds when the roll is below the capture rate.

| Capture rate | Chance |
|---|---|
| 255 | ~99.6% |
| 45 | ~17.6% |
| 3 | ~1.2% |

Some legendary Pokémon (such as Necrozma) have a capture rate of 255 because they are scripted catches in the games. To keep them challenging, the rate is capped:

- **Legendary:** at most 45 (~17.6%)
- **Mythical:** at most 30 (~11.7%)

## Caching

`internal/pokecache` is a thread-safe in-memory cache:

- Keys are request URLs and values are the response bodies as `[]byte`.
- A background goroutine removes entries older than the configured interval (10 minutes by default, set in `main.go`).
- A `sync.RWMutex` lets many reads run in parallel, while writes and cleanup get exclusive access.

## Project structure

```text
.
├── main.go          # entry point: builds the client and starts the REPL
├── repl.go          # input loop and parsing
├── commands.go      # REPL commands and catch logic
├── config.go        # state shared between commands
└── internal
    ├── pokeapi      # PokeAPI client, endpoints and response models
    └── pokecache    # in-memory cache with expiration
```

## Running tests

```sh
go test -race ./...
```

The `pokeapi` tests never reach the real API. They run against a local fake server (`net/http/httptest`) that answers with real PokeAPI responses saved in `internal/pokeapi/testdata`.

| Test file | What it covers |
|---|---|
| `repl_test.go` | Input parsing |
| `commands_test.go` | Catch rates and caps, command argument errors |
| `internal/pokecache/pokecache_test.go` | Add/Get, expiration, concurrent access |
| `internal/pokeapi/pokeapi_test.go` | Pagination, encounters, pokemon + species merge, caching, HTTP errors |

## Authorship

I wrote the application code myself while following the Boot.dev course, along with the first tests: `TestAddGet` and `TestReapLoop` for the cache, and the original `TestCleanInput` cases for the REPL.

Everything after that was written or reviewed by [Claude Code](https://claude.com/claude-code), Anthropic's AI coding assistant:

- **Code comments:** the doc comments across the project.
- **Tests:**
  - `internal/pokeapi/pokeapi_test.go`, including its fixtures in `testdata`. I started the fake server helper from an example Claude Code provided, and Claude Code completed it and wrote the tests.
  - `commands_test.go`.
  - The cache tests after `TestReapLoop`.
  - The extra `TestCleanInput` cases, plus the `continue` that keeps a length mismatch from panicking.
- **Small code changes the tests needed:**
  - Finished moving the API base URL into the `Client` struct.
  - Made `makeGetRequest` return an error for non-200 responses.
- **This README.**

Claude Code also reviewed my code during development and pointed out bugs, which I fixed myself.
