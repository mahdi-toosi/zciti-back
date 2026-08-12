# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the API server (reads config/zciti.toml)
go run cmd/main/main.go
air                       # hot reload (.air.toml, builds ./tmp/project.o)

# Build a release binary
make release              # vendored build -> ./tmp/main

# Database lifecycle (all via cmd/seeder/seeder.go flags)
make migrate              # AutoMigrate schema.MainDBModels()
make generateNecessaryData
make seed
make drop                 # drop all tables
make rebuild              # drop + wipe storage + migrate + generate + seed

# Swagger (writes docs/swagger.json)
make swag

# Tests
make tests                # gotestsum over every module's ./test/... with -p 1
go test ./app/module/order/test/ -run TestIndex_Success -v   # single test
go test ./app/module/order/test/ -v                          # single module
```

Tests run with `-p 1` on purpose: every module's `test` package shares one PostgreSQL database
(`zciti_test`), so parallel packages would clobber each other's tables.

## Configuration

`config.NewConfig()` reads `config/zciti.toml` and **panics** if it is missing or has no JWT secret.
Everything under `config/` except `example.toml` is gitignored, so a fresh clone needs
`config/zciti.toml` (server) and `config/zciti-test.toml` (tests) created from `example.toml`.
Tests load their config by absolute path via `config.ParseConfig(path, true)` — the second arg
switches from "name in ./config" lookup to a literal file path.

`docker-compose-dev.yaml` brings up Postgres (5432), a chat Postgres (5433), and Redis (6489).
The chat DB and Redis are wired in code but currently commented out.

## Architecture

Modular monolith on Fiber v2 + GORM + uber-go/fx dependency injection.

**Wiring.** `cmd/main/main.go` is the composition root: it `fx.Provide`s the infrastructure
(config, logger, Fiber app, database, payment gateways, MessageWay SMS, Bale bot, cron, middleware,
router) then lists every `<module>.Module`. `fx.Invoke(bootstrap.Start)` registers middleware +
routes and starts listening in an `fx.Hook`; the DB connects *after* the listener starts.

**A module** (`app/module/<name>/`) is always the same five packages plus an `index.go`:

| File | Contains |
|---|---|
| `index.go` | `Router` struct, `RegisterRoutes(cfg)`, and `var Module = fx.Options(...)` providing repository → service → controller → router |
| `controller/` | `index.go` aggregates controllers; `rest_controller.go` holds `IRestController` + handlers |
| `service/` | `IService` — business logic; freely imports *other modules'* services/repositories (e.g. order/service pulls coupon, product, wallet, reservation, uniwash, transaction) |
| `repository/` | `IRepository` — all GORM access, takes `*database.Database` |
| `request/` | request DTOs with `validate:` tags and `ToDomain()` |
| `response/` | response DTOs with `FromDomain()` |

Adding a module means: create the folder, add `<name>.Module` to `cmd/main/main.go`, add the
`*<name>.Router` field + constructor param + `RegisterRoutes` call in `app/router/api.go`, and add
the schema to `schema.MainDBModels()` (order matters — it drives migration order).
`app/module/example/` is a non-registered reference implementation to copy from.

**Routing and authorization.** Routes are versioned inline and almost all are business-scoped:
`/v1/business/:businessID/<resource>`. The standard guard pair is
`mdl.Protected(cfg)` (JWT) followed by `mdl.BusinessPermission(mdl.D<Domain>, mdl.P<Action>)`.
`Protected` decodes a `JWTCustomClaim` carrying a whole `schema.User` — including its
`Permissions map[businessID][]UserRole` — and puts it in `c.Locals("user")`, so permission checks
hit no database. The role → domain → action matrix lives in `app/middleware/authorization.go`;
admins bypass it. **Because permissions ride inside the token, a role change only takes effect on
re-login.**

**Responses.** Handlers return `response.Resp(c, response.Response{Data: ..., Meta: paging})`,
which emits `{Code, Messages, Data, Meta}` and substitutes `[]any{}` for nil `Data`. Errors are
returned up the stack (never written directly) and land in `response.ErrorHandler`, which maps
validator errors → 422 with translated field messages, `*fiber.Error`/`*response.Error` → their
code, everything else → 500, and mirrors each error to a Bale bot chat when configured.

**Validation.** `response.ParseAndValidate(c, req)` parses the body and runs
go-playground/validator with a translator; return its error unchanged.

**Pagination.** `paginator.Paginate(c)` reads the `page` and `itemPerPage` query params (note: not
`limit`) and returns a `*Pagination` you set on the request DTO. If `page` is absent both page and
limit become 0 — repositories treat that as "no pagination, return everything", so guard with
`if req.Pagination.Page > 0` before applying `Offset`/`Limit` and counting the total.
(The `.cursor/rules/pagination.mdc` snippet describes a different `page`/`limit` scheme that this
code does not implement — trust `utils/paginator/paginator.go`.)

## Conventions

- Method receivers are `_i` (`func (_i *service) Index(...)`); DB structs use `_db`.
- Interfaces are `I`-prefixed (`IService`, `IRepository`, `IRestController`); the implementing
  struct is unexported (`service`, `repo`, `controller`) and the constructor is the capitalized
  name returning the interface (`func Service(...) IService`).
- Named return values with a bare `return` are the norm across repositories and services.
- User-facing error messages are Persian; logs and internal errors are English.
- JSON tags are usually `json:",omitempty"` with no explicit name — the Go field name becomes the
  key, so **response DTO field names are the API contract**.
- Multi-table writes go through `repo.BeginTransaction()` and pass the `*gorm.DB` down (see
  `order/service`).
- Add Swagger annotations to handlers; regenerate with `make swag`.

## Tests

Each module has `app/module/<name>/test/` containing `setup.go` (a `TestApp` harness), optional
`mock_service.go`, and `<name>_e2e_test.go`. These are real end-to-end tests: `SetupTestApp(t)`
loads `config/zciti-test.toml`, connects to `zciti_test`, migrates, builds a real Fiber app with
the real middleware and router, and returns helpers — `CreateTestUser`, `CreateTestBusiness`,
`CreateTestOrder`, `GenerateTestToken`, `MakeRequest`, `Cleanup`/`CleanupAll`. Mirror the existing
naming (`TestIndex_Success`, `TestIndex_Unauthorized`, `TestStore_ValidationError`) and always
`defer ta.Cleanup()`. Note `mock_service.go` mocks the *service* while still using the real
repository and database — it exists to skip payment-gateway calls, not to avoid the DB.

## Vendored dependencies

`vendor/` is committed and `make release` builds with `-mod=vendor`. After changing `go.mod`, run
`go mod vendor`.

## Agent skills

### Issue tracker

Issues and specs are local markdown under `.scratch/<feature-slug>/`, one file per ticket. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical roles, unchanged, recorded as a `Status:` line in each issue file. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context — one `CONTEXT.md` and `docs/adr/` at the repo root. See `docs/agents/domain.md`.
