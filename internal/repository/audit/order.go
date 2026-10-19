package audit

import (
	"fmt"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	models.AuditOrderByObjID:     "obj_id",
	models.AuditOrderByObjDomain: "obj_domain",
	models.AuditOrderByObjName:   "obj_name",
	models.AuditOrderByActorID:   "actor_id",
	models.AuditOrderByAction:    "action",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
