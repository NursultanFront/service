// Package notificationmem - это in-memory реализация Storer для
// notificationbus.
//
// Данные хранятся в срезе внутри процесса - при перезапуске теряются, между
// несколькими инстансами сервиса не шарятся. Для текущего этапа этого
// достаточно; когда уведомлениям понадобится переживать рестарт или быть
// доступными сразу нескольким инстансам - заменить на настоящее хранилище
// (Postgres, Redis) по образцу пакетов *db/*cache в остальном проекте.
package notificationmem

import (
	"context"
	"sync"

	"github.com/ardanlabs/service/business/domain/notificationbus"
)

// Store - это in-memory хранилище уведомлений.
//
// mu (RWMutex, раз Query только читает) защищает items от одновременного
// чтения/записи из разных горутин - без него параллельные запросы словили
// бы race condition (go run/test с флагом -race это поймает). items - сам
// список уведомлений; обычного среза достаточно для текущих нужд, но мапа
// по ключу n.ID позволила бы в будущем легко добавить QueryByID без
// перебора всего среза.
type Store struct {
	mu    sync.RWMutex
	items []notificationbus.Notification
}

// NewStore constructs the in-memory store.
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
