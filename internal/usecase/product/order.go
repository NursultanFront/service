package product

import (
	productmodels "github.com/ardanlabs/service/internal/product/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProductID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the productmodels.OrderBy* constants (migration step (b2)).
const (
	OrderByProductID = productmodels.OrderByProductID
	OrderByUserID    = productmodels.OrderByUserID
	OrderByName      = productmodels.OrderByName
	OrderByCost      = productmodels.OrderByCost
	OrderByQuantity  = productmodels.OrderByQuantity
)
