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

## 2. Организация `internal/` (Два архитектурных подхода)

Внутри `internal/` находятся:
- **Стандартные инфраструктурные пакеты**: `app`, `config`, `errors`, `generated`, `server` (или внутри `handler`), `service`, `sys`, `workers`.
- **Бизнес-слои**, которые организуются одним из двух методов ниже.
  - **По умолчанию для новых и средних проектов используй Метод 1 (Разбиение по слоям / Layered Architecture)**, если пользователь явно не попросил Метод 2 (Доменное разбиение).

---

### Метод 1: Разбиение по слоям (Layered Architecture / Package by Layer)

На верхнем уровне `internal/` папки делятся **по архитектурным слоям**, а внутри каждого слоя — по доменным сущностям (`user`, `payment`, `chat` и т.д.):

```text
internal/
├── app/                        # Инициализация зависимостей (Dependency Injection руками) и запуск серверов/воркеров
│   └── app.go
├── config/                     # Чтение конфигурации из ENV / флагов запуска (только технические параметры)
│   └── config.go
├── errors/                     # Кастомные ошибки (разделение внутренней ошибки для логов и внешней для клиента)
│   └── errors.go
├── generated/                  # Автоматически сгенерированный код (Protobuf/gRPC, OpenAPI). Руками НЕ редактировать!
│
├── handler/                    # СЛОЙ 1: Транспортный слой (HTTP / gRPC точки входа)
│   ├── http/
│   │   ├── server.go           # Настройка роутера, CORS, middleware, запуск HTTP-сервера
│   │   ├── dto.go              # Входные/выходные JSON-структуры (Request / Response DTO)
│   │   └── user_handler.go     # Валидация входа, маппинг DTO <-> Models, вызов Usecase
│   └── grpc/
│       ├── server.go
│       └── user_handler.go
│
├── usecase/                    # СЛОЙ 2: Слой бизнес-логики (Сценарии использования)
│   ├── user/
│   │   ├── deps.go             # Интерфейсы внешних зависимостей для этого пакета (userRepository и др.)
│   │   └── user.go             # Структура UserUsecase и её методы (Get, Create, Update)
│   └── payment/
│       ├── deps.go
│       └── payment.go
│
├── repository/                 # СЛОЙ 3: Слой хранения данных (Работа с БД)
│   ├── user/
│   │   ├── query.go            # SQL-запросы (опционально отдельно, либо внутри user_repository.go)
│   │   └── user_repository.go  # Реализация SQL-запросов и метрик времени. 0% бизнес-логики!
│   └── payment/
│       └── payment_repository.go
│
├── models/                     # СЛОЙ 4: Чистые бизнес-модели (используются во всём проекте)
│   ├── user.go                 # struct User, struct UserShortInfo
│   └── payment.go              # struct Payment
│
├── service/                    # Клиенты к внешним сервисам (не классические БД)
│   ├── file_storage/           # Работа с S3 / MinIO (загрузка фото, документов)
│   └── docling/                # Внешние API (конвертация документов, внешние провайдеры)
│
├── sys/                        # Системно-инфраструктурные пакеты (не влияют на бизнес-логику)
│   ├── db/                     # Инициализация пула подключений к БД, менеджер транзакций
│   ├── logger/                 # Обёртка над логгером (zap / slog)
│   ├── metrics/                # Сбор системных и бизнес-метрик (Prometheus)
│   └── circuit_breaker/        # Паттерн Circuit Breaker для защиты от падения БД/внешних ресурсов
│
└── workers/                    # Фоновые периодические процессы (горутины с циклом select + time.After)
    └── task_producer.go
```

---

### Метод 2: Разбиение по доменам (Domain-Driven / Package by Feature)

Используется в крупных проектах с множеством независимых доменов. Внутри `internal/` папки верхнего уровня создаются **под каждую доменную сущность** (`user`, `payments`, `chats`, `ad`), а слои (`usecase`, `repository`, `dto`, `models`) лежат внутри каждого домена:

```text
internal/
├── app/
├── config/
├── errors/
├── generated/
├── handler/                    # Внешние HTTP и gRPC хендлеры
│   ├── grpc/
│   └── http/
├── server/                     # Настройка и запуск HTTP и gRPC серверов
│   ├── grpc/
│   └── http/
├── service/                    # Клиенты к внешним сервисам (S3/file_storage и др.)
├── sys/                        # Инфраструктура (db, logger, metrics, circuit_breaker)
├── workers/                    # Фоновые воркеры
│
├── user/                       # ДОМЕН: Пользователь
│   ├── dto/                    # Структуры передачи данных (DTO)
│   ├── models/                 # Внутренние бизнес-модели домена
│   ├── repository/             # Запросы к БД для сущности User
│   │   ├── query.go
│   │   └── user_repository.go
│   └── usecase/                # Бизнес-сценарии домена User
│       ├── deps.go
│       └── user_usecase.go
│
├── payments/                   # ДОМЕН: Платежи
│   ├── dto/
│   ├── models/
│   ├── repository/
│   └── usecase/
│
└── chats/                      # ДОМЕН: Чаты
    ├── dto/
    ├── models/
    ├── repository/
    └── usecase/
```

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
// internal/usecase/user/deps.go
package user

import (
	"context"
	"project/internal/models"
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
- **Не используй `google/wire`** (библиотека заархивирована). Собирай все зависимости явно («руками») в конструкторах внутри `internal/app/app.go`.

---

## 4. Эталонный пример связки слоёв в коде

### 1. Бизнес-модель (`internal/models/user.go`)
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

### 2. Репозиторий (`internal/repository/user/user_repository.go`)
```go
package user

import (
	"context"
	"project/internal/models"

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

### 3. Юскейс и его зависимости (`internal/usecase/user/deps.go` и `user.go`)
```go
// internal/usecase/user/deps.go
package user

import (
	"context"
	"project/internal/models"
)

type userRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
}
```

```go
// internal/usecase/user/user.go
package user

import (
	"context"
	"errors"
	"project/internal/models"
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
