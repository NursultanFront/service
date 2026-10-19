package audit

import "github.com/ardanlabs/service/internal/models"

// QueryFilter holds the available fields a query can be filtered on.
//
// NOTE(AI): alias of models.AuditQueryFilter (migration step (b2)).
type QueryFilter = models.AuditQueryFilter
