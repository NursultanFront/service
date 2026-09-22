// Package usercache contains user related CRUD functionality with caching.
package usercache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/business/domain/userbus"
	"github.com/ardanlabs/service/business/sdk/order"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/ardanlabs/service/business/types/name"
	"github.com/ardanlabs/service/business/types/role"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Store manages the set of APIs for user data and caching.
//
// TODO(dev): the actual Redis read/write/invalidate logic used to live here
// (via viccon/sturdyc, an in-process cache) and has been stripped out ahead
// of the Redis migration - see the TODO markers below in Create/Update/
// Delete/QueryByID/QueryByEmail. `cache` and `ttl` are wired up and ready to
// use; `cache` may be nil (e.g. in tests), so guard against that.
type Store struct {
	log    *logger.Logger
	storer userbus.Storer
	cache  *redis.Client
	ttl    time.Duration
	inTran bool
}

// NewStore constructs the api for data and caching access. cache may be nil,
// in which case this Store behaves as a plain passthrough with no caching.
func NewStore(log *logger.Logger, storer userbus.Storer, cache *redis.Client, ttl time.Duration) *Store {
	return &Store{
		log:    log,
		storer: storer,
		cache:  cache,
		ttl:    ttl,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (userbus.Storer, error) {
	txStorer, err := s.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log:    s.log,
		storer: txStorer,
		cache:  s.cache,
		ttl:    s.ttl,
		inTran: true,
	}

	return &store, nil
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr userbus.User) error {
	if err := s.storer.Create(ctx, usr); err != nil {
		return err
	}

	s.writeOrInvalidate(ctx, usr)

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr userbus.User) error {
	if err := s.storer.Update(ctx, usr); err != nil {
		return err
	}

	s.writeOrInvalidate(ctx, usr)

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr userbus.User) error {
	if err := s.storer.Delete(ctx, usr); err != nil {
		return err
	}

	s.deleteCache(ctx, usr)

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, page page.Page) ([]userbus.User, error) {
	return s.storer.Query(ctx, filter, orderBy, page)
}

// Count returns the total number of cards in the DB.
func (s *Store) Count(ctx context.Context, filter userbus.QueryFilter) (int, error) {
	return s.storer.Count(ctx, filter)
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (userbus.User, error) {
	if !s.inTran {
		if cachedUsr, ok := s.readCache(ctx, userID.String()); ok {
			return cachedUsr, nil
		}
	}

	usr, err := s.storer.QueryByID(ctx, userID)
	if err != nil {
		return userbus.User{}, err
	}

	s.writeOrInvalidate(ctx, usr)

	return usr, nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (userbus.User, error) {
	if !s.inTran {
		if cachedUsr, ok := s.readCache(ctx, email.Address); ok {
			return cachedUsr, nil
		}
	}

	usr, err := s.storer.QueryByEmail(ctx, email)
	if err != nil {
		return userbus.User{}, err
	}

	s.writeOrInvalidate(ctx, usr)

	return usr, nil
}

// readCache performs a safe search in the cache for the specified key.
//
// NOTE(AI): тело этой функции и toBusUserFromCache ниже написаны Claude
// (AI-ассистентом) по прямой просьбе разработчика, как разовое исключение
// из режима "AI только ревьюит" (см. AGENTS.md) - не итог самостоятельной
// работы разработчика, учитывайте при ревью.
func (s *Store) readCache(ctx context.Context, key string) (userbus.User, bool) {
	if s.cache == nil {
		return userbus.User{}, false
	}

	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		return userbus.User{}, false
	}

	var cu cacheUser
	if err := json.Unmarshal(data, &cu); err != nil {
		s.log.Error(ctx, "usercache: unmarshal", "ERROR", err)
		return userbus.User{}, false
	}

	usr, err := toBusUserFromCache(cu)
	if err != nil {
		s.log.Error(ctx, "usercache: convert", "ERROR", err)
		return userbus.User{}, false
	}

	return usr, true
}

// toBusUserFromCache converts the JSON-safe cacheUser shape back into a
// userbus.User, re-validating each primitive field through its type's
// parser rather than assigning it directly.
//
// NOTE(AI): написано Claude, см. пометку над readCache выше.
func toBusUserFromCache(cu cacheUser) (userbus.User, error) {
	id, err := uuid.Parse(cu.ID)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse id: %w", err)
	}

	nme, err := name.Parse(cu.Name)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse name: %w", err)
	}

	roles, err := role.ParseMany(cu.Roles)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse roles: %w", err)
	}

	department, err := name.ParseNull(cu.Department)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse department: %w", err)
	}

	usr := userbus.User{
		ID:           id,
		Name:         nme,
		Email:        mail.Address{Address: cu.Email},
		Roles:        roles,
		PasswordHash: cu.PasswordHash,
		Department:   department,
		Enabled:      cu.Enabled,
		DateCreated:  cu.DateCreated,
		DateUpdated:  cu.DateUpdated,
	}

	return usr, nil
}

// writeOrInvalidate populates the cache outside a transaction, but only
// invalidates the entry while inside one. A transactional write is not yet
// committed and may be rolled back, which would leave the cache holding a row
// that no longer exists in the database.
func (s *Store) writeOrInvalidate(ctx context.Context, bus userbus.User) {
	if s.inTran {
		s.deleteCache(ctx, bus)
		return
	}

	s.writeCache(ctx, bus)
}

// writeCache performs a safe write to the cache for the specified userbus.
//
// NOTE(AI): тело этой функции и toCacheUser/cacheUser ниже написаны Claude
// по прямой просьбе разработчика, разовое исключение из AGENTS.md.
func (s *Store) writeCache(ctx context.Context, bus userbus.User) {
	if s.cache == nil {
		return
	}

	data, err := json.Marshal(toCacheUser(bus))
	if err != nil {
		s.log.Error(ctx, "usercache: marshal", "ERROR", err)
		return
	}

	if err := s.cache.Set(ctx, bus.ID.String(), data, s.ttl).Err(); err != nil {
		s.log.Error(ctx, "usercache: write", "ERROR", err)
	}
}

// toCacheUser converts a userbus.User into the JSON-safe cacheUser shape,
// pulling each strong-typed field down to its primitive value.
//
// NOTE(AI): написано Claude, см. пометку над writeCache выше.
func toCacheUser(bus userbus.User) cacheUser {
	roles := make([]string, len(bus.Roles))
	for i, r := range bus.Roles {
		roles[i] = r.String()
	}

	return cacheUser{
		ID:           bus.ID.String(),
		Name:         bus.Name.String(),
		Email:        bus.Email.Address,
		Roles:        roles,
		PasswordHash: bus.PasswordHash,
		Department:   bus.Department.String(),
		Enabled:      bus.Enabled,
		DateCreated:  bus.DateCreated,
		DateUpdated:  bus.DateUpdated,
	}
}

// cacheUser is the JSON-safe shape stored in Redis. Unlike userbus.User, every
// field here is a plain primitive - name.Name and role.Role only implement
// MarshalText (for logging), not UnmarshalText, so json.Unmarshal straight
// into userbus.User would silently leave those fields zeroed instead of
// erroring.
//
// NOTE(AI): написано Claude, см. пометку над writeCache выше.
type cacheUser struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	PasswordHash []byte    `json:"password_hash"`
	Department   string    `json:"department"`
	Enabled      bool      `json:"enabled"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
}

// deleteCache performs a safe removal from the cache for the specified userbus.
func (s *Store) deleteCache(ctx context.Context, bus userbus.User) {
	if s.cache == nil {
		return
	}

	if err := s.cache.Del(ctx, bus.ID.String()).Err(); err != nil {
		s.log.Error(ctx, "usercache: delete", "ERROR", err)
	}
}
