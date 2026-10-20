package repository

import (
	"fmt"

	auditmodels "github.com/ardanlabs/service/internal/audit/models"

	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	auditmodels.OrderByObjID:     "obj_id",
	auditmodels.OrderByObjDomain: "obj_domain",
	auditmodels.OrderByObjName:   "obj_name",
	auditmodels.OrderByActorID:   "actor_id",
	auditmodels.OrderByAction:    "action",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
