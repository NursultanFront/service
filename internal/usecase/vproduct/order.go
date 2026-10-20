package vproduct

import (
	"github.com/ardanlabs/service/internal/sys/order"
	vproductmodels "github.com/ardanlabs/service/internal/vproduct/models"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProductID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the vproductmodels.OrderBy* constants (migration step (b2)).
const (
	OrderByProductID = vproductmodels.OrderByProductID
	OrderByUserID    = vproductmodels.OrderByUserID
	OrderByName      = vproductmodels.OrderByName
	OrderByCost      = vproductmodels.OrderByCost
	OrderByQuantity  = vproductmodels.OrderByQuantity
	OrderByUserName  = vproductmodels.OrderByUserName
)
