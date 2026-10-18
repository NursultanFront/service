package user

import (
	"fmt"

	"github.com/ardanlabs/service/internal/sys/order"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
)

var orderByFields = map[string]string{
	userbus.OrderByID:      "user_id",
	userbus.OrderByName:    "name",
	userbus.OrderByEmail:   "email",
	userbus.OrderByRoles:   "roles",
	userbus.OrderByEnabled: "enabled",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
