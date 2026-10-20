package usecase

import (
	"context"

	homemodels "github.com/ardanlabs/service/internal/home/models"
	usermodels "github.com/ardanlabs/service/internal/user/models"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	userbus "github.com/ardanlabs/service/internal/user/usecase"
	"github.com/google/uuid"
)

// Storer declares the behavior this package needs to persist and retrieve data.
//
// NOTE(AI): alias of homemodels.Storer (migration step (b2)): the repository's
// NewWithTx must return this exact type, so it lives where both layers can see it.
type Storer = homemodels.Storer

// userBusiness declares only the user behavior this package calls.
//
// NOTE(AI): narrowed from userbus.ExtBusiness by Claude by direct request
// (layout refactor to Method 1, deps.go on the consumer side).
type userBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (userbus.ExtBusiness, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (usermodels.User, error)
}
