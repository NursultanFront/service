package home

import (
	"fmt"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	models.HomeOrderByID:     "home_id",
	models.HomeOrderByType:   "type",
	models.HomeOrderByUserID: "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
