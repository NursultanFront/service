package usecase

import (
	homemodels "github.com/ardanlabs/service/internal/home/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the homemodels.OrderBy* constants (migration step (b2)).
const (
	OrderByID     = homemodels.OrderByID
	OrderByType   = homemodels.OrderByType
	OrderByUserID = homemodels.OrderByUserID
)
