package product

import (
	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProductID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the models.ProductOrderBy* constants (migration step (b2)).
const (
	OrderByProductID = models.ProductOrderByProductID
	OrderByUserID    = models.ProductOrderByUserID
	OrderByName      = models.ProductOrderByName
	OrderByCost      = models.ProductOrderByCost
	OrderByQuantity  = models.ProductOrderByQuantity
)
