package audit_test

import (
	"time"

	auditmodels "github.com/ardanlabs/service/internal/audit/models"

	auditapp "github.com/ardanlabs/service/internal/handler/http/audit"
)

func toAppAudit(bus auditmodels.Audit) auditapp.Audit {
	return auditapp.Audit{
		ID:        bus.ID.String(),
		ObjID:     bus.ObjID.String(),
		ObjDomain: bus.ObjDomain.String(),
		ObjName:   bus.ObjName.String(),
		ActorID:   bus.ActorID.String(),
		Action:    bus.Action,
		Data:      string(bus.Data),
		Message:   bus.Message,
		Timestamp: bus.Timestamp.Format(time.RFC3339),
	}
}

func toAppAudits(audits []auditmodels.Audit) []auditapp.Audit {
	app := make([]auditapp.Audit, len(audits))
	for i, adt := range audits {
		app[i] = toAppAudit(adt)
	}

	return app
}
