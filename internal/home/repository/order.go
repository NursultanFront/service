package repository

import (
	"fmt"

	homemodels "github.com/ardanlabs/service/internal/home/models"

	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	homemodels.OrderByID:     "home_id",
	homemodels.OrderByType:   "type",
	homemodels.OrderByUserID: "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
