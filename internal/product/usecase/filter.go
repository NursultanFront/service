package usecase

import (
	productmodels "github.com/ardanlabs/service/internal/product/models"
)

// QueryFilter holds the available fields a query can be filtered on.
//
// NOTE(AI): alias of productmodels.QueryFilter (migration step (b2)).
type QueryFilter = productmodels.QueryFilter
