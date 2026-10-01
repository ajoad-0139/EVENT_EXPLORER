# Event Explorer

Event Explorer is a web app built with Go and [Beego v2](https://github.com/beego/beego). Search for a city (with Google Places autocomplete) and browse upcoming **Music** and **Sports** events from the Ticketmaster Discovery API. You can open an event's details and jump to the ticket page.

## Features

- City search with Google Places autocomplete (cities only), proxied through the backend so the Google key stays private
- Music and Sports events for a city, fetched **concurrently** with goroutines
- Event detail page and a safe ticket redirect
- In-memory caching with TTLs (events list 5 min, single event 10 min) and `X-Cache: HIT/MISS` response headers
- Cache invalidation API
- Ticket redirect allow-list: https only, no user info, approved hosts only (blocks lookalike domains)
- Friendly error pages that map upstream failures to sensible statuses
- Swagger / OpenAPI documentation
- Unit tests for controllers, requests, routers, secrets and utils

## Tech stack

| Area | Choice |
| --- | --- |
| Language | Go 1.24+ (tests use `t.Context()`) |
| Framework | Beego v2 |
| Config | `conf/app.conf` + `.env` (via Viper) |
| Cache | Beego in-memory cache |
| External APIs | Ticketmaster Discovery API v2, Google Places API (New) |

## Project structure

```
event-explorer/
├── conf/app.conf            # port, run mode, upstream base URLs
├── controllers/
│   ├── default.go           # home page
│   ├── events.go            # list, single event, ticket redirect
│   ├── places.go            # autocomplete + place lookup (JSON)
│   ├── invalidate-cache.go  # cache invalidation (JSON)
│   └── controller_test.go
├── models/                  # error, events and location types
├── requests/                # Ticketmaster + Google clients, caches
├── routers/                 # route table (+ router_test.go)
├── secrets/                 # loads API keys and base URLs
├── static/                  # css, js, img
├── swagger/                 # OpenAPI spec + Swagger UI page
├── tests/                   # Beego sample test
├── utils/                   # error rendering, JSON helpers, ticket URL check
├── views/                   # error, event-card, index, lists, single-event .tpl
├── .env.example
├── go.mod
└── main.go
```

## Getting started

### Prerequisites

- Go 1.24 or newer
- A [Ticketmaster API key](https://developer.ticketmaster.com/)
- A Google Cloud API key with **Places API (New)** enabled



## Configuration

### `.env` (secrets, never commit this file)

```dotenv
APP_ENV=dev
GOOGLE_API_KEY=your-google-key-here
TICKETMASTER_API_KEY=your-ticketmaster-key-here
```

Real environment variables with the same names also work and are bound by Viper.

### `conf/app.conf`

```properties
appname = event-explorer
httpport = 8080
runmode = dev

googlebaseurl=https://places.googleapis.com/
ticketmasterbaseurl=https://app.ticketmaster.com/discovery/v2/
```

| Key | Purpose |
| --- | --- |
| `httpport` | Port the server listens on |
| `runmode` | `dev` or `prod`. Swagger is auto-served only in `dev` |
| `googlebaseurl` | Google Places base URL |
| `ticketmasterbaseurl` | Ticketmaster Discovery base URL (keep the trailing `/`) |

### Setup

```bash
git clone https://github.com/ajoad-0139/EVENT_EXPLORER.git
cd event-explorer
cp .env.example .env     # then edit .env with your real keys
go mod download
go run main.go
```

Open http://localhost:8080.

The app exits on startup if a key or base URL is missing, and the error lists exactly what is missing.

## Routes

### Pages (HTML)

| Method | Path | Description |
| --- | --- | --- |
| GET | `/` | Home page with the search box |
| GET | `/events?city=&countryCode=` | Music and Sports 6 events each for a city |
| GET | `/events/:eventId` | Event details |
| GET | `/redirect/:eventId` | 302 redirect to the approved ticket URL |

### API (JSON)

| Method | Path | Description |
| --- | --- | --- |
| GET | `/api/locations/autocomplete?input=&sessionToken=` | City suggestions (`input` 2-100 chars) |
| GET | `/api/locations/:placeId?sessionToken=` | Resolve a place id to `{city, countryCode}` |
| DELETE | `/api/cache/events?city=&countryCode=` | Drop a cached event list |
| DELETE | `/api/cache/events/:eventId` | Drop a cached single event |
| DELETE | `/api/cache/all` | Clear all caches |

Errors from JSON endpoints look like:

```json
{ "Error": "session token is not present" }
```

### Example requests

```bash
# autocomplete
curl "http://localhost:8080/api/locations/autocomplete?input=Toron&sessionToken=abc123"

# resolve the chosen place
curl "http://localhost:8080/api/locations/<placeId>?sessionToken=abc123"

# events page (HTML)
curl -i "http://localhost:8080/events?city=Toronto&countryCode=CA"

# clear everything
curl -X DELETE http://localhost:8080/api/cache/all
```

## API documentation (Swagger)

The OpenAPI 3 spec is in [`swagger/swagger.yml`](swagger/swagger.yml). With `runmode = dev`, Beego serves the `swagger/` folder, so after starting the app visit:

http://localhost:8080/swagger/index.html


## How it works

1. The browser calls `/api/locations/autocomplete` while the user types. The backend forwards to Google with the server-side key.
2. When a suggestion is chosen, `/api/locations/:placeId` returns the city and country code, using the same session token to keep Google billing per session.
3. `/events` checks the list cache. On a miss it calls Ticketmaster for Music and Sports in parallel, stores the result for 5 minutes, and renders `lists.tpl`.
4. `/events/:eventId` does the same with a 10 minute cache.
5. `/redirect/:eventId` validates the ticket URL against the allow-list before redirecting.

### Upstream error mapping (pages)

| Upstream status | Page status |
| --- | --- |
| 400 | 400 |
| 404 | 404 |
| 429 | 429 |
| 408 / 504 | 504 |
| anything else | 502 |

## Testing

```bash
go test ./...                              # run everything
go test ./... -coverprofile=coverage.out   # with coverage
go tool cover -html=coverage.out           # view coverage in the browser
```

Tests spin up fake Ticketmaster and Google servers with `httptest`, so no real keys or network are needed. They use `//go:linkname` to override the unexported config in the `secrets` package during tests only.

