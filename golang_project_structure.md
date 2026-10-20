# Go Project Architecture & Directory Structure Specification (LLM Instructions)

Используй этот документ как **обязательный архитектурный стандарт** при создании, рефакторинге и генерации кода в Go-проекте. Строго соблюдай разделение ответственности между папками и слоями.

---

## 1. Корневая структура репозитория (Root Layout)

В корне проекта располагаются только следующие директории и конфигурационные файлы:

```text
project-root/
├── api/                        # Внешние контракты взаимодействия (НЕ Go-код!)
│   ├── protobuf/               # .proto файлы для gRPC
│   │   └── server.proto
│   └── swagger.json            # OpenAPI / Swagger спецификации
│
├── bin/                        # Скомпилированные бинарные файлы (локальный кэш тулзов и билд сервиса)
│
├── cmd/                        # Точки входа в приложение (пакеты main)
│   ├── app/
│   │   └── main.go             # Главная точка входа (запускает основное приложение)
│   └── cron_<name>/            # (Опционально) Отдельные точки входа для периодических cron-задач
│       └── main.go
│
├── db/                         # Работа со схемой и тестовыми данными БД
│   ├── migrations/             # Статические SQL-миграции (up / down)
│   │   ├── 000001_init.up.sql
│   │   └── 000001_init.down.sql
│   └── fixtures/               # SQL-скрипты с тестовыми данными для локальной разработки
│       └── seed_test_data.sql  # Использовать INSERT ... ON CONFLICT DO NOTHING
│
├── docker/                     # Кастомные Dockerfile и инфраструктурные конфиги контейнеров
├── docs/                       # Документация проекта, гайды по локальному запуску (localrun.md), AI-инструкции
├── internal/                   # Приватный код приложения (вся бизнес-логика и внутренняя инфраструктура)
├── pkg/                        # Публичные Go-пакеты (Shared-код, доступный для импорта другими репозиториями)
├── scripts/                    # Вспомогательные shell/bash скрипты (.sh) для генерации, CI/CD и автоматизации
├── tools/                      # Вспомогательные Go-утилиты разработки (mockery, golangci-lint, migrate, gofumpt)
│   ├── go.mod
│   └── tools.go                # //go:generate директивы для сборки утилит в папку /bin
│
├── .env.example                # Шаблон переменных окружения
├── docker-compose.yaml         # Локальная инфраструктура (PostgreSQL, Redis, MinIO/S3 и т.д.)
├── Makefile                    # Команды запуска, миграций, генерации и линтинга
├── go.mod
└── go.sum
```

### Строгие правила для корневых директорий:
1. **`cmd/`**: В `main.go` **запрещено** писать бизнес-логику. `main.go` должен только инициализировать контекст и вызывать сборку приложения из `internal/app`. Главный `main` в сервисе всегда один.
2. **`pkg/` vs `api/`**:
   - В `pkg/` кладётся **только публичный Go-код** (общие структуры событий Kafka, общие клиенты/контракты на Go), который может импортироваться другими микросервисами.
   - В `api/` кладутся **только языково-независимые контракты** (`.proto`, `swagger`/`openapi`).
3. **`tools/` и `bin/`**: Любые внешние генераторы и линтеры (`mockery`, `golangci-lint`, `migrate`, `gofumpt`) фиксируются в `tools/` и компилируются в `bin/`, чтобы не зависеть от глобального окружения разработчика.
4. **`db/fixtures/`**: Для локального дебага создаются идемпотентные SQL-скрипты (`ON CONFLICT DO NOTHING`), чтобы не наполнять базу вручную после каждого перезапуска.

---

## 2. Организация `internal/` (Доменное разбиение — единственный обязательный вариант)

С 2026-10-20 проект использует **Разбиение по доменам (Domain-Driven Layout / Package by Domain)** по образцу рабочего проекта ментора Германа (`typespace`). Разбиение по слоям на верхнем уровне (`internal/usecase`, `internal/repository`, `internal/models`) **не используется** — не создавай эти папки. История перехода: `DOMAIN_LAYOUT_PLAN.md`.

Внутри `internal/` находятся:
- **Инфраструктурные пакеты** (общие для всех доменов): `app`, `config`, `errors`, `generated`, `handler`, `server`, `service`, `sys`, `types`, `workers`.
- **Доменные папки** (`user`, `product`, `home`, `audit`, `vproduct`, …): каждая владеет своими слоями `models`, `usecase`, `repository`.

```text
internal/
├── app/                        # Сборка зависимостей руками (DI) и запуск серверов/воркеров: app/<svc>/
├── config/                     # Конфигурация из ENV / флагов (только технические параметры): config/<svc>/
├── errors/                     # Ошибки: внутренняя причина (для логов) + безопасное сообщение (для клиента)
├── generated/                  # Сгенерированный код (Protobuf/gRPC, OpenAPI). Руками НЕ редактировать!
│
├── handler/                    # СЛОЙ 1: Транспорт (HTTP / gRPC хендлеры)
│   ├── http/
│   │   ├── mid/                # HTTP middleware
│   │   └── user/               # user_handler.go, dto.go, filter.go, order.go, route.go
│   └── grpc/
│       └── auth/               # auth_handler.go
│
├── server/                     # Настройка и запуск серверов (роутер, CORS, middleware, lifecycle)
│   ├── http/
│   │   └── server.go
│   └── grpc/
│       └── server.go
│
├── service/                    # Клиенты к внешним сервисам (S3, другие API, authclient)
├── sys/                        # Инфраструктура без бизнес-логики (db, metrics, auth, circuit_breaker)
├── types/                      # Общие строгие типы (Value Objects): name, money, role, password, ...
├── workers/                    # Фоновые периодические процессы
│
├── user/                       # 📦 ДОМЕН: Пользователь (всё про User в одной папке)
│   ├── models/                 # СЛОЙ 4: бизнес-модели домена
│   │   ├── user.go             # User, NewUser, UpdateUser, Storer, ErrNotFound, ErrUniqueEmail
│   │   └── query.go            # QueryFilter, OrderBy* константы
│   ├── usecase/                # СЛОЙ 2: бизнес-логика домена
│   │   ├── deps.go             # Интерфейсы зависимостей (на стороне потребителя)
│   │   ├── user.go             # Business / ExtBusiness
│   │   └── extensions/         # Декораторы (otel, audit, ...)
│   └── repository/             # СЛОЙ 3: работа с БД (0% бизнес-логики)
│       ├── user_repository.go
│       ├── dto.go              # DB-структуры и конвертеры toDB<Type> / toBus<Type>
│       ├── filter.go
│       └── order.go
│
└── product/                    # 📦 ДОМЕН: Продукты (та же форма: models/, usecase/, repository/)
    ├── models/
    ├── usecase/
    └── repository/
```

### Правила доменной раскладки:
1. **Имена пакетов совпадают с папками** (`package models`, `package usecase`, `package repository`). При импорте используй алиасы: `<domain>models`, `<domain>bus` (usecase), `<domain>db` (repository) — например `usermodels`, `userbus`, `userdb`.
2. **Типы внутри домена без префикса домена**: `usermodels.Storer`, `usermodels.QueryFilter`, `usermodels.OrderByEmail`, `usermodels.ErrNotFound` (не `UserStorer`, `UserQueryFilter`). Имена сущностей (`User`, `NewUser`, `Product`) остаются.
3. **`repository` импортирует только `internal/<domain>/models`** (и `internal/types/*`), никогда `usecase`.
4. **Общие строгие типы — в `internal/types/*`**, чтобы домены переиспользовали их без циклических зависимостей. Домен не импортирует `models` другого домена без необходимости; если нужно — только `models`, не `usecase`/`repository`.
5. **Серверы — в `internal/server/{http,grpc}`**, хендлеры — в `internal/handler/{http,grpc}`.

---

## 3. Архитектурные правила и паттерны кода (Обязательно к соблюдению)

### Правило 1: Строгое разделение ответственности между 3 слоями (`Handler` -> `Usecase` -> `Repository`)

1. **`handler` (Транспорт и первичная валидация)**:
   - Принимает запрос по HTTP или gRPC.
   - Выполняет **только первичную техническую валидацию**: проверка на пустые поля, проверка `id > 0`, ограничение `limit` (например, запрет запрашивать 1 000 000 записей, максимум 300), проверка корректности формата.
   - Парсит входной JSON/Protobuf в `DTO`, конвертирует (мапит) его в бизнес-модель `models` и передаёт в `usecase`.
   - Получает результат от `usecase`, конвертирует в выходной `DTO` и отдаёт клиенту.

2. **`usecase` (Чистая бизнес-логика)**:
   - Вся бизнес-логика проекта живёт **исключительно здесь**.
   - Примеры бизнес-логики: проверка `if user.Deactivated`, `if user.Banned`, расчёт комиссий, проверка прав доступа, оркестрация вызовов нескольких репозиториев или внешних сервисов (`service`).
   - **Не дробить файлы без необходимости**: все стандартные сценарии одной сущности (`Get`, `Create`, `Update`, `Delete`) объединяй в одну структуру `UserUsecase`, а не создавай отдельную папку под каждую функцию.

3. **`repository` (Чистая работа с БД)**:
   - Содержит **только** SQL-запросы, сериализацию/десериализацию данных для БД и замер времени выполнения запроса (метрики).
   - **Категорически запрещено** писать бизнес-логику и бизнес-условия внутри `repository`. Что репозиторию передали — то он и записал/прочитал.

---

### Правило 2: Инверсия зависимостей (Dependency Inversion Principle — `deps.go`)

- **Интерфейсы всегда объявляются в том пакете, где они используются (на стороне потребителя), а не там, где они реализуются.**
- Пакет `usecase` **не должен** напрямую зависеть от конкретной структуры пакета `repository` или конкретного внешнего провайдера (например, `CloudPayments`).
- В папке каждого юскейса создавай файл `deps.go`, в котором описывай нужные этому юскейсу интерфейсы:

```go
// internal/user/usecase/deps.go
package usecase

import (
	"context"
	models "project/internal/user/models"
)

// Интерфейс объявляется у потребителя (в пакете usecase)
type userRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
}

type paymentProvider interface {
	ProcessPayment(ctx context.Context, userID string, amount int64) error
}
```

---

### Правило 3: Разделение `models` и `dto`

- **`models`** — внутренние бизнес-сущности системы (`User`, `Payment`).
  - Никогда не храни и не передавай сырой пароль, используй `PasswordHash`.
  - Если для разных сценариев требуется отдавать урезанный набор полей (например, по требованиям безопасности/комплаенса нельзя возвращать чувствительные поля профиля), создавай отдельную структуру `UserShortInfo` или мапь только разрешённые поля в выходной `DTO`.
- **`dto` (Data Transfer Object)** — структуры с тегами `json:"..."` или `db:"..."`, которые используются на границах системы (в `handler` для приёма/отдачи по сети или в `repository` для чтения специфичных структур БД).

---

### Правило 4: Обработка ошибок (`internal/errors`)

- Никогда не отдавай наружу пользователю сырые ошибки базы данных или инфраструктуры (например, `pq: connection refused`).
- Используй обогащённые ошибки:
  - **Внутренняя ошибка (`err`)** — пишется в серверные логи (`sys/logger`).
  - **Внешнее сообщение (`message`)** — безопасный, понятный человеку текст, который отдаётся клиенту в HTTP/gRPC ответе (например: *«Сервис временно недоступен, попробуйте позднее»* или *«Пользователь уже зарегистрирован»*).

---

### Правило 5: Отказ от пакета-свалки `utils` и кодогенераторов DI (`wire`)

- **Не создавай безликий пакет `utils/`** для всего подряд (со временем он превращается в нечитаемую свалку). Если функция относится к конкретному домену — держи её рядом с доменом; если это системный хелпер — клади в осмысленный пакет внутри `sys/` или `pkg/`.
- **Не используй `google/wire`** (библиотека заархивирована). Собирай все зависимости явно («руками») в конструкторах внутри `internal/app/<svc>/app.go`.

---

## 4. Эталонный пример связки слоёв в коде

### 1. Бизнес-модель (`internal/user/models/user.go`)
```go
package models

import "time"

type User struct {
	ID           string
	Name         string
	Surname      string
	Email        string
	PasswordHash string
	Deactivated  bool
	Banned       bool
	CreatedAt    time.Time
}

type UserShortInfo struct {
	ID      string
	Name    string
	Surname string
}
```

### 2. Репозиторий (`internal/user/repository/user_repository.go`)
```go
package repository

import (
	"context"
	models "project/internal/user/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*models.User, error) {
	const query = `
		SELECT id, name, surname, email, password_hash, deactivated, banned, created_at
		FROM users
		WHERE id = $1
	`
	var u models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Surname, &u.Email,
		&u.PasswordHash, &u.Deactivated, &u.Banned, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
```

### 3. Юскейс и его зависимости (`internal/user/usecase/deps.go` и `user.go`)
```go
// internal/user/usecase/deps.go
package usecase

import (
	"context"
	models "project/internal/user/models"
)

type userRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
}
```

```go
// internal/user/usecase/user.go
package usecase

import (
	"context"
	"errors"
	models "project/internal/user/models"
)

var (
	ErrUserDeactivated = errors.New("user is deactivated")
	ErrUserBanned      = errors.New("user is banned")
)

type Usecase struct {
	repo userRepository
}

func NewUsecase(repo userRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetPublicProfile(ctx context.Context, id string) (*models.UserShortInfo, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Бизнес-правила проверяются только в Usecase:
	if user.Deactivated {
		return nil, ErrUserDeactivated
	}
	if user.Banned {
		return nil, ErrUserBanned
	}

	return &models.UserShortInfo{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
	}, nil
}
```

### 4. Сборка зависимостей (`internal/app/app.go` и `cmd/app/main.go`)
```go
// cmd/app/main.go
package main

import (
	"context"
	"log"
	"project/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
```
