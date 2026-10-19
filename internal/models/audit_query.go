package models

import (
	"time"

	"github.com/ardanlabs/service/internal/models/types/domain"
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/audit/{filter,order}.go in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// AuditQueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type AuditQueryFilter struct {
	ObjID     *uuid.UUID
	ObjDomain *domain.Domain
	ObjName   *name.Name
	ActorID   *uuid.UUID
	Action    *string
	Since     *time.Time
	Until     *time.Time
}

// Set of fields that audit results can be ordered by.
const (
	AuditOrderByObjID     = "a"
	AuditOrderByObjDomain = "b"
	AuditOrderByObjName   = "c"
	AuditOrderByActorID   = "d"
	AuditOrderByAction    = "e"
)
