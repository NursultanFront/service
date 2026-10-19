package user

import (
	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the models.UserOrderBy* constants (migration step (b2)).
const (
	OrderByID      = models.UserOrderByID
	OrderByName    = models.UserOrderByName
	OrderByEmail   = models.UserOrderByEmail
	OrderByRoles   = models.UserOrderByRoles
	OrderByEnabled = models.UserOrderByEnabled
)
