package audit

import auditbus "github.com/ardanlabs/service/internal/audit/usecase"

var orderByFields = map[string]string{
	"obj_id":     auditbus.OrderByObjID,
	"obj_domain": auditbus.OrderByObjDomain,
	"obj_name":   auditbus.OrderByObjName,
	"actor_id":   auditbus.OrderByActorID,
	"action":     auditbus.OrderByAction,
}
