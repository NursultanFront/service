package user

import (
	usermodels "github.com/ardanlabs/service/internal/user/models"
)

// QueryFilter holds the available fields a query can be filtered on.
//
// NOTE(AI): alias of usermodels.QueryFilter (migration step (b2)).
type QueryFilter = usermodels.QueryFilter
