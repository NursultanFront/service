package usecase

import (
	auditmodels "github.com/ardanlabs/service/internal/audit/models"
)

// QueryFilter holds the available fields a query can be filtered on.
//
// NOTE(AI): alias of auditmodels.QueryFilter (migration step (b2)).
type QueryFilter = auditmodels.QueryFilter
