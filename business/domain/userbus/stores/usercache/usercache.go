// Package usercache contains user related CRUD functionality with caching.
package usercache

import (
	"context"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/business/domain/userbus"
	"github.com/ardanlabs/service/business/sdk/order"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/business/sdk/sqldb"
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
// TODO(dev): implement using s.cache (a *redis.Client). Remember values need
// to be serialized (redis stores bytes/strings, not Go structs) - e.g. with
// encoding/json. Guard against s.cache == nil.
func (s *Store) readCache(ctx context.Context, key string) (userbus.User, bool) {
	return userbus.User{}, false
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
// TODO(dev): implement using s.cache and s.ttl (redis SET with expiration -
// see the redis package's SetEx/Set(...).WithTTL). Guard against s.cache ==
// nil. The old sturdyc version keyed by both ID and email; consider whether
// you want the same two keys, or a single key with the other as a lookup.
func (s *Store) writeCache(ctx context.Context, bus userbus.User) {
}

// deleteCache performs a safe removal from the cache for the specified userbus.
//
// TODO(dev): implement using s.cache (redis DEL). Guard against s.cache ==
// nil.
func (s *Store) deleteCache(ctx context.Context, bus userbus.User) {
}
