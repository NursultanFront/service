package user

import "github.com/ardanlabs/service/internal/sys/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID      = "a"
	OrderByName    = "b"
	OrderByEmail   = "c"
	OrderByRoles   = "d"
	OrderByEnabled = "e"
)
