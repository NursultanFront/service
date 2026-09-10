// Package notificationbus provides a business logic layer for handling
// notifications, typically created from Kafka events.
//
// TODO(dev): это каркас. Нужно:
//  1. Написать реализацию Storer (например, notificationbus/stores/notificationdb
//     для Postgres, или стор на Redis) и подключить её в NewBusiness.
//  2. Дописать Create/Query ниже — сейчас они ничего не делают, просто
//     возвращают нулевые значения.
package notificationbus

import (
	"context"
	"slices"

	"github.com/ardanlabs/service/foundation/logger"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, n Notification) error
	Query(ctx context.Context) ([]Notification, error)
}

// ExtBusiness interface provides support for extensions that wrap extra
// functionality around the core business logic.
type ExtBusiness interface {
	Create(ctx context.Context, nn NewNotification) (Notification, error)
	Query(ctx context.Context) ([]Notification, error)
}

// Extension is a function that wraps a new layer of business logic around
// the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for notification access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a notification business API for use.
func NewBusiness(log *logger.Logger, storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		log:    log,
		storer: storer,
	})

	for _, ext := range slices.Backward(extensions) {
		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

// Create adds a new notification to the system.
//
// TODO(dev): открой business/domain/auditbus/auditbus.go и посмотри на его
// Create (около 60-й строки) — здесь должно быть по той же схеме:
//  1. Собери значение Notification из nn: сгенерируй новый ID через
//     uuid.New() (импорт "github.com/google/uuid"), перенеси ProductID/
//     EventType/Message из nn, выставь CreatedAt через time.Now() (импорт
//     "time"). Присвой всё это локальной переменной, например
//     `n := Notification{...}`.
//  2. Вызови b.storer.Create(ctx, n). Проверь ошибку — если не nil, верни
//     Notification{}, fmt.Errorf("create notification: %w", err) (импорт
//     "fmt") — такое оборачивание ошибки — общий паттерн по всему проекту,
//     оно оставляет исходную ошибку доступной через errors.Is/errors.As.
//  3. При успехе верни n, nil.
func (b *Business) Create(ctx context.Context, nn NewNotification) (Notification, error) {
	return Notification{}, nil
}

// Query retrieves the list of existing notifications.
//
// TODO(dev): этот метод проще, чем Create:
//  1. Вызови ns, err := b.storer.Query(ctx).
//  2. Если err != nil — верни nil, fmt.Errorf("query notifications: %w", err).
//  3. Иначе верни ns, nil.
//
// Когда это заработает end-to-end, вернись и подумай: нужна ли тебе
// поддержка фильтра/сортировки/пагинации, как у остальных bus-пакетов
// (глянь сигнатуру auditbus.Query — она принимает QueryFilter, order.By,
// page.Page)? Для первой рабочей версии это не обязательно — списки
// уведомлений, скорее всего, небольшие — но полезно знать, что такой
// паттерн в проекте есть, если фича вырастет.
func (b *Business) Query(ctx context.Context) ([]Notification, error) {
	return nil, nil
}
