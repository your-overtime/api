# API v3 – Design (draft)

Status: proposal for discussion. Nothing here is implemented yet.

## Goals

- Serve everything the new web UI and the trackers (browser extensions,
  desktop apps, Android) need from one API.
- Usage events from the trackers move into the API (today they live in the
  web UI's own SQLite database).
- **v2 stays backward compatible**: same paths, request and response formats
  and auth methods. v2 and v3 work on the same data, and both report the same
  numbers.

## Non-goals (for v3.0)

- Teams, organisations, roles, billing.
- Projects, clients and tags (can follow in a v3.x without breaking changes).

## Architecture

```
api repo
├── v2/          frozen; stays buildable for reference and rollback
└── v3/          new Go module, one binary
    ├── migrations/      versioned SQL migrations
    ├── internal/
    │   ├── store/       database access (MariaDB, SQLite)
    │   ├── calc/        work day calculation, pure functions
    │   ├── service/     business logic, used by both API versions
    │   ├── apiv3/       /api/v3 handlers
    │   └── apiv2/       /api/v2 compatibility layer on top of service/
    └── main.go
```

- **One binary serves `/api/v2` and `/api/v3`.** Both use the same service
  layer, so there is exactly one calculation and one set of rules. The v2
  handlers are ported, not rewritten: same routes, same JSON field names
  (PascalCase), same status codes.
- **Shared database.** v3 adopts the existing v2 tables (`employees`,
  `tokens`, `activities`, `holidays`, `webhooks`) unchanged and only adds
  tables. Rolling back to the `:2` image stays possible (see Migration).
- **Versioned SQL migrations** replace GORM `AutoMigrate`. Migration 1
  describes the existing v2 schema and is a no-op on existing databases.
- MariaDB stays the production database; SQLite for tests and small setups.

## Data model

Existing (unchanged): users (`employees`), tokens, activities, holidays,
webhooks.

New:

| Table | Purpose |
|---|---|
| `work_schedules` | Working time per user, valid from a date until the next schedule: week minutes, working days (fixed or dynamic count), holidays per year. Replaces the settings on the user. |
| `balance_adjustments` | Manual corrections of the overtime balance on a day (e.g. carry-over from a former system), with a note. Replaces writing work days. |
| `usage_events` | Spans of app or browser tab use uploaded by trackers. Moved from the web UI. |
| `schema_migrations` | Applied migrations. |

`work_days` (the former cache) is no longer used. It stays in the database
for the comparison tool and for a rollback to v2.

### Work day calculation

Days are calculated on every request from activities, holidays, schedules and
adjustments (one query each, in memory; a whole year takes milliseconds).
Nothing is cached, so changes to the past are always reflected. The
calculation is ready on the branch `feature/work-schedules` (`calc.go` with
tests) and moves to `v3/internal/calc`.

That branch also fixes three bugs of the v2 calculation, which change
balances once when v3 goes live:

- dynamic working days: worked days were never counted as working days, so
  a full week showed as overtime,
- activities across midnight were not counted,
- the first week of a year was estimated instead of calculated.

### Usage events

```
id        client-chosen, unique per user (idempotent uploads)
user_id
source    "chrome", "firefox", "macos", "linux", "android", ...
device    free text, e.g. host name
start, end
app, title, url   url without query and fragment
location  e.g. "Büro", derived by the client; empty = unknown
```

- Same validation as today: max 1000 events per batch, max 24 h per event,
  no events in the future, app or url required.
- Events never change the balance. They are shown in the timeline and turned
  into activity suggestions that the user confirms.
- Privacy: they belong to the user only (no access for other users or
  admins via the API); per-user retention (default e.g. 90 days, older
  events are deleted); option to store no titles/URLs.

## API v3

Conventions

- JSON with camelCase fields, timestamps RFC 3339, dates `YYYY-MM-DD`.
- Errors: `{"error": {"code": "invalid_schedule", "message": "..."}}`.
- Auth: `Authorization: Bearer <token>` (also accepts `Token <token>`).
  No tokens in query parameters, except the iCal feed URLs.
- Token scopes: `read`, `write`, `events` (trackers only need `events`).
  Existing tokens map to `read` (readonly) or `read write events`.

Endpoints (sketch)

| Method | Path | |
|---|---|---|
| POST | `/auth/token` | Basic auth → new token (login for apps and the web UI) |
| GET | `/me` | user and the schedule valid today |
| PATCH | `/me` | name, login, password |
| GET | `/overview?date=` | as v2, plus balance since the first day |
| GET | `/days?from=&to=` | calculated days: active, target, holiday, adjustment, overtime |
| GET POST PATCH DELETE | `/activities` | `GET` with `from`, `to`, `q` |
| GET | `/timer` | running activity |
| POST | `/timer/start`, `/timer/stop` | |
| GET POST PATCH DELETE | `/holidays` | |
| GET POST DELETE | `/schedules` | `POST` with the same `validFrom` replaces |
| GET POST DELETE | `/adjustments` | |
| GET POST | `/events` | `POST` takes the current tracker payload unchanged |
| GET | `/stream` | Server-Sent Events, see below |
| GET POST DELETE | `/tokens`, `/webhooks` | |
| GET | `/activities.ics`, `/holidays.ics` | |

Live updates (`GET /stream`, SSE): `activity.changed`, `timer.started`,
`timer.stopped`, `holiday.changed`, `schedule.changed`, `events.added`.
Payload is the affected ID; clients reload what they show. Every client
(web, desktop, Android) sees a timer started elsewhere at once.

## v2 compatibility

| v2 | Behaviour in v3 |
|---|---|
| all paths, fields, status codes | unchanged |
| `GET /overview`, `/workday` | calculated by v3 (same numbers as v3) |
| `PATCH /account` with working time fields | creates a schedule valid from today; `GET /account` shows the schedule valid today |
| `POST /workday` (set a day's overtime) | stored as balance adjustment: requested overtime minus calculated overtime of that day |
| `POST /user?adminToken=` | unchanged, also creates the first schedule |

Ensured by tests: the existing v2 API tests run against the compatibility
layer, plus golden-file tests for every v2 response.

## Migration and rollout

1. Deploy `:3` (serves v2 and v3). On start it runs the migrations: every
   user gets a schedule with the current settings, valid since 1970-01-01.
2. Before that, run `overtime -compare-workdays` against a copy of the
   production database: it lists per user the days where the former cache
   and the new calculation differ. Add schedules for former settings where
   needed; inform users whose balance changes because of the bug fixes.
3. Trackers switch their upload URL to `/api/v3/events`; the web UI
   imports the events stored so far (idempotent by ID) and then works
   against v3 only.
4. Rollback: `:2` keeps working on the same database, because v3 only adds
   tables. Data written only to new tables (schedules, adjustments, events)
   is ignored by v2; its cache must be cleared (`DELETE FROM work_days`).

## Phases

1. **Skeleton**: v3 module, migrations, adopt v2 schema, v2 compatibility
   layer with the existing v2 tests green.
2. **Calculation**: schedules, adjustments, comparison tool; v2 overview on
   the new calculation.
3. **API v3**: endpoints above, token scopes, OpenAPI docs.
4. **Events and live updates**: `/events`, `/stream`, retention; trackers
   and web UI switch.

## Open questions

1. camelCase in v3 (consistent with JS clients) or PascalCase like v2?
2. Keep the global admin token, or admin users with a scope?
3. Default retention for usage events, and should titles/URLs be optional?
4. Should the balance carry over between years (today the overview shows
   the current year only)?
5. Projects and tags in v3.0 or later?
