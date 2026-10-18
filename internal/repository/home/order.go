package home

import (
	"fmt"

	"github.com/ardanlabs/service/internal/sys/order"
	homebus "github.com/ardanlabs/service/internal/usecase/home"
)

var orderByFields = map[string]string{
	homebus.OrderByID:     "home_id",
	homebus.OrderByType:   "type",
	homebus.OrderByUserID: "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
