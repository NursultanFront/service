// Package mid provides app level middleware support.
package mid

import (
	"context"
	"errors"

	homemodels "github.com/ardanlabs/service/internal/home/models"
	productmodels "github.com/ardanlabs/service/internal/product/models"
	usermodels "github.com/ardanlabs/service/internal/user/models"

	"github.com/ardanlabs/service/internal/sys/auth"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/pkg/foundation/web"
	"github.com/google/uuid"
)

func checkIsError(e web.Encoder) error {
	err, hasError := e.(error)
	if hasError {
		return err
	}

	return nil
}

// =============================================================================

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
	productKey
	homeKey
	trKey
)

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

// GetSubjectID returns the subject id from the claims.
func GetSubjectID(ctx context.Context) uuid.UUID {
	v := GetClaims(ctx)

	subjectID, err := uuid.Parse(v.Subject)
	if err != nil {
		return uuid.UUID{}
	}

	return subjectID
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}

func setUser(ctx context.Context, usr usermodels.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (usermodels.User, error) {
	v, ok := ctx.Value(userKey).(usermodels.User)
	if !ok {
		return usermodels.User{}, errors.New("user not found in context")
	}

	return v, nil
}

func setProduct(ctx context.Context, prd productmodels.Product) context.Context {
	return context.WithValue(ctx, productKey, prd)
}

// GetProduct returns the product from the context.
func GetProduct(ctx context.Context) (productmodels.Product, error) {
	v, ok := ctx.Value(productKey).(productmodels.Product)
	if !ok {
		return productmodels.Product{}, errors.New("product not found in context")
	}

	return v, nil
}

func setHome(ctx context.Context, hme homemodels.Home) context.Context {
	return context.WithValue(ctx, homeKey, hme)
}

// GetHome returns the home from the context.
func GetHome(ctx context.Context) (homemodels.Home, error) {
	v, ok := ctx.Value(homeKey).(homemodels.Home)
	if !ok {
		return homemodels.Home{}, errors.New("home not found in context")
	}

	return v, nil
}

func setTran(ctx context.Context, tx sqldb.CommitRollbacker) context.Context {
	return context.WithValue(ctx, trKey, tx)
}

// GetTran retrieves the value that can manage a transaction.
func GetTran(ctx context.Context) (sqldb.CommitRollbacker, error) {
	v, ok := ctx.Value(trKey).(sqldb.CommitRollbacker)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return v, nil
}
