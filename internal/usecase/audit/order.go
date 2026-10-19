package audit

import (
	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByObjID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the models.AuditOrderBy* constants (migration step (b2)).
const (
	OrderByObjID     = models.AuditOrderByObjID
	OrderByObjDomain = models.AuditOrderByObjDomain
	OrderByObjName   = models.AuditOrderByObjName
	OrderByActorID   = models.AuditOrderByActorID
	OrderByAction    = models.AuditOrderByAction
)
