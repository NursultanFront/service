# План миграции к стандарту `golang_project_structure.md`

Стандарт структуры — [`golang_project_structure.md`](../golang_project_structure.md)
(Метод 1, разбиение по слоям). Правила для агентов и порядок старшинства документов —
[`AGENTS.md`](../AGENTS.md), раздел *Project Structure & Architecture Standard*.
Карта старых путей Ardan Labs → текущие — [`zarf/prompts/project-structure.md`](../zarf/prompts/project-structure.md).

Принципы миграции:

- **Один шаг — одна ветка/PR.** Шаги не смешиваются с фичами и друг с другом.
- **Без смены логики**, кроме шага (h) (ошибки), где поведение меняется осознанно.
- Паттерны кода внутри слоёв (строгие типы `internal/models/types/*`, конвертеры
  `toBus*`/`toDB*`/`fromBus*Response`/`toApp*`, декораторы `ExtBusiness`/`Extension`)
  **сохраняются** — меняется только структура.
- Проверка после каждого шага минимум: `go vet ./...` и `make test-only`; для шагов,
  которые трогают сборку/запуск, ещё `make compose-build-up` (и `docker image prune -f`).
- Статус на 2026-10-18: ветка `refactor/go-standard-layout` уже частично перестроена
  (`internal/{handler,usecase,repository,models,service,sys,errors,generated}`,
  `api/protobuf`, `pkg/foundation`), изменения ещё не закоммичены.

---

## Текущее состояние vs стандарт

| Область | Сейчас (на диске) | Стандарт | Шаг |
|---|---|---|---|
| Сервисы в репо | `cmd/sales`, `cmd/auth`, `cmd/metrics` + CLI `cmd/admin`, `cmd/logfmt` | один сервис на репозиторий, `cmd/app/main.go` | (a), (i) |
| `main.go` | толстые: `cmd/sales/main.go` ≈340 строк, `cmd/auth/main.go` ≈320, конфиг + wiring + запуск | тонкий `main`, вызывает `internal/app` | (a) |
| Wiring / DI | `cmd/<svc>/main.go` + `cmd/{sales,auth}/build/*` (build-теги `crud`/`reporting`) | `internal/app/app.go`, руками | (a) |
| Конфиг | `conf.Parse` внутри `main.go` | `internal/config/config.go` | (a) |
| Транспорт | `internal/handler/http/<x>/`, `internal/handler/grpc/auth/` | `internal/handler/{http,grpc}/` | ✅ (имена файлов — (b)) |
| Бизнес-логика | `internal/usecase/<x>/` (+ `deps.go` уже есть) | `internal/usecase/<entity>/deps.go + <entity>.go` | ✅ (имена файлов — (b)) |
| БД-слой | `internal/repository/<x>/` | `internal/repository/<entity>/<entity>_repository.go` | ✅ (имена файлов — (b)) |
| Модели | `internal/models/<x>.go`, `internal/models/types/*` | `internal/models/` | ✅ |
| Имена файлов | `<x>bus.go`, `<x>db.go`, `<x>app.go`, `model.go`, `mux.go` | `<x>.go`, `<x>_repository.go`, `<x>_handler.go`, `dto.go`, `server.go` | (b) |
| Ошибки | `internal/errors` (`errs.go`, `codes.go`): код + сообщение | внутренняя причина (лог) + внешнее сообщение (клиент) | (h) |
| Инфраструктура | `internal/sys/*` + `pkg/foundation/*` (logger, web, otel, keystore, worker, docker) | всё внутреннее — `internal/sys/*`; `pkg/` только публичное | (f) |
| Контракты | `api/protobuf/authapp.proto` ✅, но в `api/` лежит фронтенд `api/frontends/admin` | `api/` — только `.proto`/OpenAPI | (g) |
| Деплой/инфра | `zarf/{docker,compose,k8s,helm,keys,prompts}` | `docker/`, `docker-compose.yaml` + `.env.example` в корне, `docs/`, `scripts/` | (c) |
| Миграции БД | darwin, один файл `internal/sys/db/migrate/sql/migrate.sql` + `seed.sql` (embed) | `db/migrations/NNNNNN_*.up/down.sql`, `db/fixtures/*.sql` | (d) |
| Dev-тулзы | `go install …@latest` в `makefile` (staticcheck, govulncheck, goimports, protoc-gen-go…) | `tools/` (свой `go.mod`, `tools.go`) → `bin/` | (e) |
| Фоновые процессы | нет `internal/workers` | `internal/workers/` при появлении | — |

### ✅ gRPC-фиксы авторизации — перенесены 2026-10-18

> Выполнено: перенесены в `internal/handler/grpc/auth/authapp.go` и
> `internal/service/authclient/grpc/model.go`; проверено временным in-process тестом
> (Token без/с неверными creds → Unauthenticated, верные → токен с ролью USER; Authenticate
> отдаёт Claims; rule_any/user_only/admin_or_subject проходят, rule_admin_only → PermissionDenied),
> на старом коде тест падает. Ниже — исходное описание.

Фиксы сделаны **только** в worktree-ветке `claude/keen-proskuriakova-20f1cb` на старой
раскладке (`app/domain/grpcauthapp/authapp.go`, `app/sdk/authclient/grpc/model.go`,
не закоммичены). В `refactor/go-standard-layout` их нет (проверено 2026-10-18:
нет `ndToInt64`, `int64ToND` не nil-safe). Перенести в новые пути **до шага (b)**,
чтобы не переносить дважды:

| Фикс | Куда |
|---|---|
| `Authenticate` заполняет `Claims` в ответе | `internal/handler/grpc/auth/authapp.go` |
| nil-safe `ndToInt64` / `int64ToND` | `internal/service/authclient/grpc/model.go` |
| interceptor: проверка по `Auth_Token_FullMethodName` только для `Token`; `Token` берёт claims из `mid.GetClaims` | `internal/handler/grpc/auth/authapp.go` |
| отказ OPA → `codes.PermissionDenied` | `internal/handler/grpc/auth/authapp.go` |

Проверка: `go vet ./...`, `make test-only`, ручной вызов gRPC `Authenticate`/`Token` через
`make compose-build-up` с `SALES_AUTH_*` в режиме `grpc`.

---

## Шаги

Порядок по умолчанию: **(0) → (a) → (b) → (f) → (h) → (e) → (c) → (d) → (g) → (i)**.
(c), (d), (e), (g) мало связаны с Go-кодом и могут идти в другом порядке, но до (i).

### (0) Зафиксировать текущую перестройку

- **Цель:** закоммитить уже сделанный перенос в `internal/*` как отдельный коммит, чтобы
  дальше двигаться маленькими шагами.
- **Что:** текущие незакоммиченные изменения ветки `refactor/go-standard-layout` +
  перенос gRPC-фиксов (см. выше).
- **Проверка:** `go vet ./...`, `make test-only`, `make compose-build-up`.
- **Зависит от:** —.

### (a) `internal/config` + `internal/app` для каждого сервиса, тонкий `cmd/*/main.go`

- **Цель:** правило 1 стандарта — `main` только создаёт контекст и вызывает `app.Run(ctx)`.
- **Что переносится:**
  - структура `cfg` + `conf.Parse` из `cmd/<svc>/main.go` → `internal/config`;
  - сборка зависимостей (DB, Redis, authclient, usecase'ы + extensions, mux/роуты,
    gRPC-сервер, debug-сервер, graceful shutdown) → `internal/app`;
  - `cmd/{sales,auth}/build/*` (выбор набора роутов по build-тегам) → `internal/app`.
- **Пока сервисы в одном репо** — по пакету на сервис, например `internal/app/sales`,
  `internal/app/auth`, `internal/app/metrics` и `internal/config/{sales,auth,metrics}`
  (после шага (i) станет просто `internal/app` и `internal/config`). Точное именование —
  открытый вопрос Q1.
- **Затрагивает:** `cmd/sales/main.go`, `cmd/auth/main.go`, `cmd/metrics/main.go`,
  `cmd/{sales,auth}/build/*`, новые `internal/app/*`, `internal/config/*`; скилл
  `business-layer-extensions` (текст уже указывает на `internal/app` как цель).
- **Риски:** build-теги `crud`/`reporting` в `cmd/sales/build` — их надо сохранить
  (или заменить конфигом — Q2); порядок shutdown; `cmd/sales/tests/*` могут зависеть от
  `build` — проверить импорты; `логирование `GOMAXPROCS` в `main`.
- **Проверка:** `go vet ./...`, `make test-only`, `make compose-build-up`, запуск через
  `go run` по AGENTS.md (*Local Development Workflow*).
- **Зависит от:** (0).

### (b) Переименование файлов под стандарт (без смены логики)

- **Цель:** имена файлов из стандарта; чисто механический шаг.
- **Что:**
  - `internal/usecase/<x>/<x>bus.go` → `<x>.go`;
  - `internal/repository/<x>/<x>db.go` → `<x>_repository.go` (`productpg.go`,
    `productsqlite.go`, `usercache.go` — Q3), `model.go` → `dto.go`;
  - `internal/handler/http/<x>/<x>app.go` → `<x>_handler.go`, `model.go` → `dto.go`;
  - `internal/handler/http/server/mux.go` → `server.go`;
  - `internal/handler/grpc/auth/authapp.go` → `auth_handler.go`, `start.go` → `server.go`;
  - тесты рядом (`<x>bus_test.go` → `<x>_test.go` и т.д.).
- **Затрагивает:** только имена файлов в `internal/{usecase,repository,handler}`; импорты
  не меняются (пакеты те же).
- **Риски:** минимальные; `git mv`, чтобы сохранить историю; не смешивать с правками кода.
  Имена пакетов (`package userbus` и т.п.) — отдельный вопрос Q4.
- **Проверка:** `go vet ./...`, `make test-only`.
- **Зависит от:** (0); лучше после (a), чтобы не конфликтовать.

### (c) `zarf/` → `docker/`, `docs/`, `scripts/`; compose и `.env.example` в корень

- **Цель:** корень репо как в стандарте.
- **Что:**
  - `zarf/docker/dockerfile.*` → `docker/`;
  - `zarf/compose/docker_compose.yaml` → `docker-compose.yaml` (корень), `pg_hba.conf` →
    `docker/`; **`zarf/compose/database-data/` — данные Postgres, не переносить в git**,
    добавить в `.gitignore` / вынести в volume;
  - создать `.env.example` со всеми `SALES_*`, `AUTH_*`, `METRICS_*` переменными;
  - `zarf/prompts/*` → `docs/` (AI-инструкции), `localrun.md` в `docs/`;
  - `zarf/k8s`, `zarf/helm` → Q5, `zarf/keys` → Q6;
  - shell-скрипты из `makefile` → `scripts/`.
- **Затрагивает:** `makefile` (≈33 упоминания `zarf`), `zarf/k8s/dev/Tiltfile`,
  Dockerfile'ы (`COPY`-пути), helm-чарты, AGENTS.md (*Local Development Workflow*
  упоминает `docker_compose.yaml`), `README.md`.
- **Риски:** сломанные пути в `makefile`/Tilt/Helm; ключи `zarf/keys/*.pem`
  монтируются в контейнеры; `database-data` случайно попадёт в git.
- **Проверка:** `make compose-build-up`, k8s/Tilt-таргеты `makefile`, `go vet ./...`.
- **Зависит от:** — (лучше после (a)).

### (d) Миграции: darwin single-file → `db/migrations` up/down + `db/fixtures`

- **Цель:** раздел про `db/` стандарта.
- **Что:** разрезать `internal/sys/db/migrate/sql/migrate.sql` (darwin, версии внутри
  одного файла) на `db/migrations/000001_<name>.up.sql` / `.down.sql`; `seed.sql` →
  `db/fixtures/seed_test_data.sql` с `INSERT … ON CONFLICT DO NOTHING`.
- **Инструмент:** предложение — `golang-migrate` (CLI через `tools/` + библиотека с
  `embed`/`iofs` для `cmd/admin migrate` и `dbtest`). Выбор — Q7.
- **Затрагивает:** `internal/sys/db/migrate`, `internal/sys/dbtest` (поднимает схему в
  тестах), `cmd/admin/commands/{migrate,seed}.go`, `makefile`, k8s init-контейнеры,
  `go.mod` (удалить `github.com/ardanlabs/darwin/v3`).
- **Риски:** **существующие БД** знают версии darwin (своя служебная таблица) — нужен
  baseline (`migrate force N`) либо пересоздание dev-БД; down-миграции пишутся впервые;
  `dbtest` должен работать с новым механизмом.
- **Проверка:** чистая БД: up → down → up; `make test-only`; `make compose-build-up`.
- **Зависит от:** (c) (папка `db/` в корне), (e) (CLI `migrate` в `bin/`).

### (e) `tools/` + `bin/`

- **Цель:** правило 3 стандарта — тулзы зафиксированы версиями.
- **Что:** `tools/go.mod`, `tools/tools.go` (blank-imports или `tool`-директивы Go 1.24+)
  для staticcheck/golangci-lint (Q8), govulncheck, goimports/gofumpt, protoc-gen-go,
  protoc-gen-go-grpc, migrate, expvarmon, hey; таргет `make tools` собирает в `bin/`;
  таргеты `makefile` вызывают `./bin/<tool>`; `bin/` в `.gitignore`.
- **Затрагивает:** `makefile` (секция `go install …@latest`, `lint`, `vuln-check`,
  proto-генерация), CI (если есть), `.gitignore`.
- **Риски:** Windows (`.exe` в `bin/`); расхождение версий с текущими глобальными.
- **Проверка:** `make tools && make lint && make vuln-check`, `make test`.
- **Зависит от:** —.

### (f) `pkg/foundation/*` → `internal/sys/*`

- **Цель:** правило 2 стандарта — в `pkg/` только публичное.
- **Что:** `logger`, `web`, `otel`, `keystore`, `worker`, `docker` → `internal/sys/<name>`;
  обновить ≈70 Go-файлов с импортом `pkg/foundation/...`; удалить `pkg/foundation`.
  В `pkg/` оставить только то, что реально понадобится другим сервисам после (i)
  (кандидаты: публичный клиент auth, будущие общие события) — решается в (i).
- **Затрагивает:** все слои (импорты), скилл `layered-architecture-types` (упоминания
  `pkg/foundation` → `internal/sys`), AGENTS.md (ссылка на `pkg/foundation/web/web.go`
  в *Local Development Workflow*).
- **Риски:** коллизии имён в `internal/sys`; `pkg/foundation/docker` — тестовый хелпер,
  возможно место в `internal/sys/dbtest`.
- **Проверка:** `go vet ./...`, `make test-only`, `rg "pkg/foundation"` → пусто.
- **Зависит от:** (0). Должен быть **до** (i).

### (g) `api/frontends/admin` вон из `api/`

- **Цель:** `api/` — только контракты.
- **Варианты (Q9 — решение за пользователем):** отдельный репозиторий (логично, если
  делаем (i)) или `frontend/admin` в корне этого репо.
- **Затрагивает:** `api/frontends/admin`, `makefile`/Dockerfile/compose, если фронт
  собирается или раздаётся оттуда; CORS (`http://localhost:3001`) не меняется.
- **Риски:** `dist/` и `node_modules` в git; пути в CI.
- **Проверка:** фронт собирается и логинится на `http://localhost:3001` (httpOnly cookie).
- **Зависит от:** —.

### (h) `internal/errors`: внутренняя причина + внешнее сообщение

- **Цель:** правило 4 стандарта.
- **Что:** тип ошибки хранит `cause error` (для логов, `errors.Unwrap`) и безопасное
  `message` для клиента; `mid/errors.go` логирует причину и отдаёт только сообщение;
  gRPC-хендлер мапит в `status.Error(code, message)`; аудит мест, где в ответ уходит
  `err.Error()` из БД/инфраструктуры.
- **Затрагивает:** `internal/errors/{errs.go,codes.go}`, `internal/handler/http/mid/errors.go`,
  все `internal/handler/http/*`, `internal/handler/grpc/auth`, тесты `cmd/sales/tests/*`.
- **Риски:** **меняется контракт ответов** — фронтенд и тесты, сравнивающие текст ошибки;
  `FieldErrors` остаются детальными (это безопасные ошибки валидации).
- **Проверка:** `make test-only` (с обновлёнными ожиданиями), ручная проверка 4xx/5xx.
- **Зависит от:** (b).

### (i) Разделение на отдельные репозитории: sales, auth, metrics

Последний и самый большой шаг. Каждый репозиторий — полный стандарт: `cmd/app/main.go`,
`internal/app`, `internal/config`, свои `db/`, `docker/`, `tools/`, `makefile`,
`docker-compose.yaml`.

- **Что куда (предложение, требует решений):**

  | Что | Репозиторий |
  |---|---|
  | sales API (`internal/handler/http/{user,product,home,vproduct,audit,tran,raw,check}`, usecase'ы, repository, models, миграции) | `sales` |
  | auth (`internal/handler/http/auth`, `oauth`, `internal/handler/grpc/auth`, `internal/sys/auth` + rego) | `auth` |
  | metrics sidecar (`cmd/metrics/*`) | `metrics` |
  | `cmd/admin` (migrate, seed, gentoken, useradd) | Q10 |
  | `cmd/logfmt` | Q11 |
  | `api/protobuf/authapp.proto` | владелец — `auth`; доставка в `sales` — Q12 |
  | `internal/generated/authpb` | генерируется у потребителя из `.proto` или публикуется (см. варианты) |
  | `internal/service/authclient` | `sales` (потребитель) — или публичный клиент `auth/pkg/authclient` |

- **Общая БД:** сейчас auth работает с таблицей `users` через `internal/usecase/user`.
  Q13: auth владеет пользователями (своя БД + API для sales) или общая БД временно
  остаётся. Это главный архитектурный риск шага.
- **Варианты шаринга кода:**
  1. **Публичный shared-модуль** (`pkg/logger`, `pkg/web`, `pkg/otel`, `pkg/authpb`… в
     отдельном репо/модуле, версии через git-теги): один источник правды; минусы —
     релизный цикл, breaking changes бьют по всем, риск превратить его в `utils`.
  2. **Копирование** `internal/sys/*` в каждый сервис: полная независимость, проще для
     обучения; минусы — дрейф, дублирование фиксов.
  3. **Гибрид** (на обсуждение): инфраструктуру (`internal/sys/*`) копировать, а
     **контракт** (`authapp.proto` + сгенерированный Go-клиент) публиковать из `auth`
     (`auth/pkg/...` или отдельный модуль контрактов) — по стандарту `pkg/` как раз для
     этого.
- **Затрагивает:** всё; `go.mod` (`module github.com/ardanlabs/service` → новые module
  path, Q14), k8s/helm, compose (сеть между сервисами), `makefile`, CI.
- **Риски:** общая БД/схема; тесты `cmd/sales/tests/*` поднимают auth в процессе
  (`internal/handler/http/apitest`) — после разделения нужен мок или контейнер; история
  git (Q15).
- **Проверка:** в каждом репо `go vet ./...`, `make test-only`; e2e — compose с тремя
  образами: логин через auth (HTTP и gRPC), CRUD в sales, метрики в metrics.
- **Зависит от:** (a), (b), (c), (d), (e), (f), (h); (g) желательно.

---

## Открытые вопросы (пользователю / ментору Герману)

1. **Q1.** Пока сервисы в одном репо: `internal/app/<svc>` + `internal/config/<svc>` — ок,
   или сначала разделить репо, а потом делать (a)?
2. **Q2.** Build-теги `crud`/`reporting` в `cmd/sales/build` — сохраняем (в `internal/app`),
   заменяем конфигом или удаляем?
3. **Q3.** Имена файлов при нескольких реализациях репозитория (`postgres/`, `sqlite/`,
   `cache/`): `product_repository.go` в каждой папке или `postgres.go`/`sqlite.go`?
4. **Q4.** Имена пакетов: оставить `package userbus`/`userdb`/`userapp` или
   переименовать в `package user` (как в примерах стандарта)?
5. **Q5.** Куда `zarf/k8s` и `zarf/helm` — стандарт не описывает (`deploy/`?).
6. **Q6.** `zarf/keys/*.pem` (dev-ключи JWT) — где хранить и коммитить ли вообще?
7. **Q7.** Инструмент миграций: `golang-migrate` (предложение), `goose` или darwin с
   разбивкой файлов? Как делаем baseline для существующих dev-БД?
8. **Q8.** Линтер: оставляем `staticcheck` или переходим на `golangci-lint` (упомянут в
   стандарте)?
9. **Q9.** `api/frontends/admin` — отдельный репозиторий или `frontend/` в этом репо?
10. **Q10.** `cmd/admin` — куда после разделения (в `sales`, в `auth`, отдельный CLI)?
11. **Q11.** `cmd/logfmt` — оставить, вынести или удалить?
12. **Q12.** Как `sales` получает `authapp.proto`: submodule, buf, копия, публичный Go-модуль?
13. **Q13.** Пользователи и БД: auth владеет `users` отдельно или общая БД пока остаётся?
14. **Q14.** Новые module path — уходим от `github.com/ardanlabs/service`?
15. **Q15.** Разделение репо с историей (`git filter-repo`) или чистый старт?
16. **Q16.** `internal/sys/auth` импортирует `internal/usecase/user`, а `internal/sys/dbtest`
    — все usecase/repository. Это противоречит «sys без бизнес-логики» — переносим
    auth-правила в `internal/usecase/auth` (+ `deps.go`), а `dbtest` — ближе к `internal/app`?
17. **Q17.** Метод 1 + декораторы `ExtBusiness` + строгие типы — ментор подтверждает, что
    это не противоречит стандарту (AGENTS.md трактует именно так)?
