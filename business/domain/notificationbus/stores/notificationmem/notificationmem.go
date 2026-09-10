// Package notificationmem is a placeholder Storer for notificationbus.
//
// TODO(dev): это заглушка — Create/Query сейчас ничего не делают. Нужно
// либо реализовать здесь реальное in-memory хранилище (слайс/мапа под
// мьютексом), либо заменить весь пакет на настоящее хранилище (Postgres,
// Redis) по образцу пакетов *db/*cache в остальном проекте.
package notificationmem

import (
	"context"
	"sync"

	"github.com/ardanlabs/service/business/domain/notificationbus"
)

// Store is a placeholder in-memory store for notifications.
//
// TODO(dev): добавь сюда два поля:
//  1. mu sync.Mutex (или sync.RWMutex, раз Query только читает) — защищает
//     слайс ниже. Без него параллельные запросы, читающие/пишущие
//     одновременно, словят race condition (go run/test с флагом -race это
//     поймает, если забудешь).
//  2. items []notificationbus.Notification — само хранилище. Обычного
//     слайса достаточно для заглушки; мапа по ключу n.ID позволила бы
//     потом легко добавить QueryByID без перебора всего слайса.
type Store struct{
	mu    sync.RWMutex
	items []notificationbus.Notification
}

// NewStore constructs the placeholder store.
//
// TODO(dev): когда добавишь поле items выше — инициализируй его здесь,
// например items: make([]notificationbus.Notification, 0). Строго не
// обязательно (append прекрасно работает и с nil-слайсом), но так яснее.
func NewStore() *Store {
	return &Store{items: make([]notificationbus.Notification, 0)}
}

// Create adds a new notification to the in-memory store.
func (s *Store) Create(ctx context.Context, n notificationbus.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = append(s.items, n)

	return nil
}

// Query returns all the notifications currently in the store.
func (s *Store) Query(ctx context.Context) ([]notificationbus.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy of the items slice to prevent external modification.
	copiedItems := make([]notificationbus.Notification, len(s.items))
	copy(copiedItems, s.items)

	return copiedItems, nil
}
