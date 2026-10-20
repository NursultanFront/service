# План перевода `internal/` на Доменную структуру (как в `typespace` у Германа)

Этот документ задаёт финальную целевую структуру `internal/` по образцу рабочего проекта ментора Германа (`typespace` — **Разбиение по доменам / Domain-Driven Layout**) и пошаговый план перехода с текущей слоистой раскладки.

> **Статус (2026-10-20): все 3 шага выполнены** — `f7374179` (шаг 1), `519e9dd9` (шаг 2), коммит с документацией (шаг 3). Решения: gRPC-сервер вынесен в `internal/server/grpc`; в моделях — алиасы импорта (`usermodels`, `userbus`, `userdb`) и типы без префикса домена (`usermodels.Storer`, `QueryFilter`, `ErrNotFound`).

---

## 1. Что остаётся без изменений (уже соответствует `typespace` на 100%)

Корень репозитория и инфраструктурные пакеты внутри `internal/` **не трогаем** — они уже приведены к эталону:
- Корень: `api/`, `bin/`, `cmd/`, `db/{migrations,fixtures}`, `deploy/`, `docker/`, `docs/`, `frontend/`, `pkg/foundation/`, `tools/`, `.env.example`, `docker-compose.yaml`, `makefile`.
- Внутри `internal/`: `app/`, `config/`, `errors/`, `generated/`, `service/`, `sys/`.

---

## 2. Целевая структура `internal/` (Доменное разбиение как у Германа)

Вместо того чтобы размазывать одну бизнес-сущность по трём разным папкам верхнего уровня (`internal/usecase/user`, `internal/repository/user`, `internal/models/user.go`), **каждый бизнес-домен собирается в свою собственную папку верхнего уровня внутри `internal/`** (без лишней прослойки `internal/pkg/`, как рекомендовал Герман на созвоне):

```text
internal/
├── app/                        # Сборка приложения (DI руками): sales/, auth/, metrics/
├── config/                     # Конфигурация из ENV: sales/, auth/, metrics/
├── errors/                     # Кастомные ошибки (Cause для логов + безопасный Message наружу)
├── generated/                  # Сгенерированный Protobuf/gRPC код (authpb/)
│
├── handler/                    # Внешние точки входа (HTTP и gRPC хендлеры)
│   ├── grpc/
│   │   └── auth/               # auth_handler.go
│   └── http/
│       ├── mid/                # HTTP middleware
│       ├── query/              # Обёртка ответов со списком и пагинацией
│       ├── apitest/            # Хелперы интеграционных тестов
│       ├── user/               # user_handler.go, dto.go, filter.go, order.go, route.go
│       ├── product/            # product_handler.go, dto.go, filter.go, order.go, route.go
│       ├── home/               # home_handler.go, dto.go, filter.go, order.go, route.go
│       ├── audit/              # audit_handler.go, dto.go, filter.go, order.go, route.go
│       ├── vproduct/           # vproduct_handler.go, dto.go, filter.go, order.go, route.go
│       ├── auth/               # auth_handler.go, dto.go, route.go
│       ├── check/              # check_handler.go, dto.go, route.go
│       ├── tran/               # tran_handler.go, dto.go, route.go
│       ├── oauth/              # oauth_handler.go, route.go
│       └── raw/                # raw_handler.go, route.go
│
├── server/                     # Настройка и запуск HTTP / gRPC серверов (на уровне internal/, как в typespace)
│   ├── grpc/
│   │   └── server.go           # бывший internal/handler/grpc/auth/server.go
│   └── http/
│       └── server.go           # бывший internal/handler/http/server/server.go
│
├── service/                    # Клиенты к внешним сервисам
│   └── authclient/             # http/ и grpc/
│
├── sys/                        # Системная инфраструктура без бизнес-логики
│   ├── auth/                   # JWT + OPA/Rego (с локальным интерфейсом UserLookup)
│   ├── db/                     # sqldb, tran, dialect, dbarray, migrate
│   ├── metrics/                # метрики + collector/ и publisher/
│   ├── rediscache/
│   ├── delegate/
│   ├── order/
│   ├── page/
│   ├── debug/
│   ├── dbtest/
│   └── unittest/
│
├── types/                      # Общие валидируемые строгие типы (Value Objects), используемые всеми доменами
│   ├── domain/
│   ├── home/
│   ├── money/
│   ├── name/
│   ├── password/
│   ├── quantity/
│   └── role/
│
├── user/                       # 📦 ДОМЕН 1: Пользователи (всё про User в одной папке)
│   ├── models/                 # Бизнес-модели, ошибки и фильтры домена User
│   │   ├── user.go             # struct User, NewUser, UpdateUser, Storer, ErrNotFound, ErrUniqueEmail
│   │   └── query.go            # QueryFilter, OrderBy константы
│   ├── repository/             # Слой БД для User (0% бизнес-логики)
│   │   ├── user_repository.go
│   │   ├── dto.go              # DB-структура dbUser и конвертеры toDBUser / toBusUser
│   │   ├── filter.go
│   │   ├── order.go
│   │   └── cache/
│   │       └── user_repository.go
│   └── usecase/                # Бизнес-сценарии User
│       ├── deps.go             # Интерфейс Storer (алиас/объявление на стороне потребителя)
│       ├── user.go             # Business / ExtBusiness
│       ├── user_test.go
│       ├── filter.go
│       ├── order.go
│       ├── event.go
│       ├── testutil.go
│       └── extensions/
│           ├── useraudit/
│           └── userotel/
│
├── product/                    # 📦 ДОМЕН 2: Продукты
│   ├── models/
│   │   ├── product.go
│   │   └── query.go
│   ├── repository/
│   │   ├── common/             # dto.go, filter.go, order.go
│   │   ├── postgres/           # product_repository.go
│   │   └── sqlite/             # product_repository.go
│   └── usecase/
│       ├── deps.go
│       ├── product.go
│       ├── product_test.go
│       ├── filter.go
│       ├── order.go
│       ├── event.go
│       ├── testutil.go
│       └── extensions/
│           └── productotel/
│
├── home/                       # 📦 ДОМЕН 3: Дома
│   ├── models/
│   │   ├── home.go
│   │   └── query.go
│   ├── repository/
│   │   ├── home_repository.go
│   │   ├── dto.go
│   │   ├── filter.go
│   │   └── order.go
│   └── usecase/
│       ├── deps.go
│       ├── home.go
│       ├── home_test.go
│       ├── filter.go
│       ├── order.go
│       ├── event.go
│       ├── testutil.go
│       └── extensions/
│           └── homeotel/
│
├── audit/                      # 📦 ДОМЕН 4: Аудит
│   ├── models/
│   │   └── audit.go
│   ├── repository/
│   │   ├── audit_repository.go
│   │   ├── dto.go
│   │   ├── filter.go
│   │   └── order.go
│   └── usecase/
│       ├── deps.go
│       ├── audit.go
│       ├── audit_test.go
│       ├── filter.go
│       ├── order.go
│       ├── testutil.go
│       └── extensions/
│           └── auditotel/
│
└── vproduct/                   # 📦 ДОМЕН 5: Витрина продуктов (Read Model)
    ├── models/
    │   └── vproduct.go
    ├── repository/
    │   ├── vproduct_repository.go
    │   ├── dto.go
    │   ├── filter.go
    │   └── order.go
    └── usecase/
        ├── deps.go
        ├── vproduct.go
        ├── vproduct_test.go
        ├── filter.go
        ├── order.go
        └── extensions/
            └── vproductotel/
```

---

## 3. Пошаговый план миграции для Claude

Выполнить в 3 чистых атомарных коммита (используя `git mv`, без изменения бизнес-логики):

### ✅ Шаг 1: Вынести общие строгие типы в `internal/types/` и разложить `internal/models/` по доменам
1. Переместить `internal/models/types/*` → `internal/types/*` (это общие примитивы/Value Objects `name`, `money`, `role`, `password`, `quantity`, `home`, `domain`, которые переиспользуются между `user`, `product`, `home`, `audit`, `vproduct` без циклических зависимостей).
2. Разложить доменные модели из `internal/models/` по соответствующим доменным папкам `internal/<domain>/models/`:
   - `internal/models/user.go` + `internal/models/user_query.go` → `internal/user/models/`
   - `internal/models/product.go` + `internal/models/product_query.go` → `internal/product/models/`
   - `internal/models/home.go` + `internal/models/home_query.go` → `internal/home/models/`
   - `internal/models/audit.go` → `internal/audit/models/`
   - `internal/models/vproduct.go` → `internal/vproduct/models/`
   - Общие доменные ошибки из `internal/models/errors.go` (`ErrNotFound`, `ErrUniqueEmail`, `ErrAuthenticationFailure`) разнести по `internal/<domain>/models/` (или оставить в `internal/errors`), чтобы каждый домен владел своими моделями и ошибками, а `repository` импортировал только `internal/<domain>/models` (не импортируя `usecase`).

### ✅ Шаг 2: Переместить `usecase` и `repository` внутрь каждого домена `internal/<domain>/` и вынести `internal/server/`
1. Для каждого из 5 доменов (`user`, `product`, `home`, `audit`, `vproduct`):
   - `git mv internal/usecase/<domain> internal/<domain>/usecase`
   - `git mv internal/repository/<domain> internal/<domain>/repository`
2. Удалить опустевшие папки верхнего уровня `internal/usecase`, `internal/repository`, `internal/models`.
3. Вынести серверные пакеты на верхний уровень `internal/server/` (как в `typespace`):
   - `git mv internal/handler/http/server internal/server/http`
   - `git mv internal/handler/grpc/auth/server.go internal/server/grpc/server.go` (или оставить gRPC-сервер в `internal/server/grpc`).
4. Обновить все пути импортов в проекте, проверить `go vet ./...` и `make test-only`.

### ✅ Шаг 3: Обновить документацию (`golang_project_structure.md`, `AGENTS.md`, `docs/MIGRATION_PLAN.md`, скиллы в `.agents/skills/`)
1. В `AGENTS.md` и `golang_project_structure.md` зафиксировать **Доменную структуру (Domain-Driven Layout как в `typespace`)** как единственную обязательную структуру проекта.
2. Обновить пути в `.agents/skills/{layered-architecture-types,business-layer-extensions,review-pr}` на `internal/<domain>/{models,usecase,repository}` и `internal/types/*`.
3. Прогнать `go vet ./...`, `make lint` и `make test-only`.
