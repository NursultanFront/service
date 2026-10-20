package usecase

import (
	auditmodels "github.com/ardanlabs/service/internal/audit/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByObjID, order.ASC)

// Set of fields that the results can be ordered by.
//
// NOTE(AI): aliases of the auditmodels.OrderBy* constants (migration step (b2)).
const (
	OrderByObjID     = auditmodels.OrderByObjID
	OrderByObjDomain = auditmodels.OrderByObjDomain
	OrderByObjName   = auditmodels.OrderByObjName
	OrderByActorID   = auditmodels.OrderByActorID
	OrderByAction    = auditmodels.OrderByAction
)
