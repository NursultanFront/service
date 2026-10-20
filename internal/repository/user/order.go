package user

import (
	"fmt"

	usermodels "github.com/ardanlabs/service/internal/user/models"

	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	usermodels.OrderByID:      "user_id",
	usermodels.OrderByName:    "name",
	usermodels.OrderByEmail:   "email",
	usermodels.OrderByRoles:   "roles",
	usermodels.OrderByEnabled: "enabled",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
