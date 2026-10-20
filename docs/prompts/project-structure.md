# Project Structure — DEPRECATED

> **This document is obsolete.** It used to describe the Ardan Labs template layout
> (`internal/app` + `internal/business` + `pkg/foundation`). That layout is no longer the
> target.
>
> - **Standard (source of truth for structure):** [`golang_project_structure.md`](../../golang_project_structure.md)
>   — Method 1, layered: `handler` → `usecase` → `repository` + `models`.
> - **Rules for agents and how the standard combines with code patterns:**
>   [`AGENTS.md`](../../AGENTS.md), section *Project Structure & Architecture Standard*.
> - **Step-by-step migration plan:** [`docs/MIGRATION_PLAN.md`](../../docs/MIGRATION_PLAN.md).
>
> Do not use this file as a blueprint for new code.

## Old path → current path

Use this map only to read old notes, commits, or prompts. "Current" = where the code is on
disk today; "target" = where the standard wants it after the migration plan is done.

### Business layer → usecase / repository / models

| Old (Ardan Labs) | Current | Notes |
|---|---|---|
| `internal/business/domain/<x>bus/` | `internal/usecase/<x>/` | `<x>` = audit, home, product, user, vproduct. Core file `<x>bus.go` → `<x>.go` |
| `internal/business/domain/<x>bus/extensions/*` | `internal/usecase/<x>/extensions/*` | ExtBusiness/Extension decorators |
| `internal/business/domain/<x>bus/stores/<x>db/` | `internal/repository/<x>/` | `<x>db.go` → `<x>_repository.go`, `model.go` → `dto.go` |
| `internal/business/domain/productbus/stores/productdb/` | `internal/repository/product/{common,postgres,sqlite}/` | |
| `internal/business/domain/userbus/stores/usercache/` | `internal/repository/user/cache/` | |
| `internal/business/domain/<x>bus/model.go` | `internal/models/<x>.go` | |
| `internal/business/types/*` | `internal/models/types/*` | domain, home, money, name, password, quantity, role |
| `internal/business/sdk/sqldb` | `internal/sys/db` | incl. `dbarray`, `dialect` |
| `internal/business/sdk/migrate` | `internal/sys/db/migrate` (golang-migrate runner) | SQL in `db/migrations` + `db/fixtures`, embedded by package `db` |
| `internal/business/sdk/dbtest` | `internal/sys/dbtest` | |
| `internal/business/sdk/order` | `internal/sys/order` | |
| `internal/business/sdk/page` | `internal/sys/page` | |
| `internal/business/sdk/delegate` | `internal/sys/delegate` | |
| `internal/business/sdk/rediscache` | `internal/sys/rediscache` | |
| `internal/business/sdk/unittest` | `internal/sys/unittest` | |

### App layer → handler / service / sys / errors

| Old (Ardan Labs) | Current | Notes |
|---|---|---|
| `internal/app/domain/<x>app/` | `internal/handler/http/<x>/` | audit, auth, check, home, oauth, product, raw, tran, user, vproduct. `<x>app.go` → `<x>_handler.go`, `model.go` → `dto.go` |
| `internal/app/domain/grpcauthapp/` | `internal/handler/grpc/auth/` | generated code → `internal/generated/authpb/`, `.proto` → `api/protobuf/authapp.proto` |
| `internal/app/sdk/mid` | `internal/handler/http/mid` | |
| `internal/app/sdk/mux` | `internal/handler/http/server` | `mux.go` → `server.go` |
| `internal/app/sdk/apitest` | `internal/handler/http/apitest` | |
| `internal/app/sdk/query` | `internal/handler/http/query` | |
| `internal/app/sdk/errs` | `internal/errors` | target: internal cause + safe external message |
| `internal/app/sdk/auth` | `internal/sys/auth` | incl. `rego/` |
| `internal/app/sdk/authclient` | `internal/service/authclient` | `http/`, `grpc/` |
| `internal/app/sdk/debug` | `internal/sys/debug` | |
| `internal/app/sdk/metrics` | `internal/sys/metrics` | |

### Foundation, entry points, deployment

| Old (Ardan Labs) | Current | Target |
|---|---|---|
| `foundation/*` | `pkg/foundation/*` | stays `pkg/foundation/*` — shared by all `cmd/*` (decision 2026-10-18) |
| `api/services/<svc>` | `cmd/<svc>` | stays `cmd/<svc>` in this repo (decision 2026-10-18) + `internal/app/<svc>` + `internal/config/<svc>` |
| `api/tooling/{admin,logfmt}` | `cmd/admin`, `cmd/logfmt` | see `docs/MIGRATION_PLAN.md`, step (i) |
| `api/frontends/admin` | `frontend/admin` | |
| `zarf/docker/*` | `docker/*` | |
| `zarf/compose/docker_compose.yaml` | `docker-compose.yaml` (root), `docker/pg_hba.conf`, data in `docker/database-data` (gitignored) | |
| `zarf/k8s`, `zarf/helm`, `zarf/keys` | `deploy/k8s`, `deploy/helm`, `deploy/keys` | |
| `zarf/prompts` | `docs/prompts` | |
