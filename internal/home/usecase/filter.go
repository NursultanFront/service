package usecase

import (
	homemodels "github.com/ardanlabs/service/internal/home/models"
)

// QueryFilter holds the available fields a query can be filtered on.
//
// NOTE(AI): alias of homemodels.QueryFilter (migration step (b2)).
type QueryFilter = homemodels.QueryFilter
