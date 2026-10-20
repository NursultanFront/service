package usecase

import (
	"github.com/ardanlabs/service/internal/sys/order"
	usermodels "github.com/ardanlabs/service/internal/user/models"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the usermodels.OrderBy* constants (migration step (b2)).
const (
	OrderByID      = usermodels.OrderByID
	OrderByName    = usermodels.OrderByName
	OrderByEmail   = usermodels.OrderByEmail
	OrderByRoles   = usermodels.OrderByRoles
	OrderByEnabled = usermodels.OrderByEnabled
)
