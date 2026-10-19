package user

import (
	"fmt"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	models.UserOrderByID:      "user_id",
	models.UserOrderByName:    "name",
	models.UserOrderByEmail:   "email",
	models.UserOrderByRoles:   "roles",
	models.UserOrderByEnabled: "enabled",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
