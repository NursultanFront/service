package vproduct

import (
	"fmt"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	models.VProductOrderByProductID: "product_id",
	models.VProductOrderByUserID:    "user_id",
	models.VProductOrderByName:      "name",
	models.VProductOrderByCost:      "cost",
	models.VProductOrderByQuantity:  "quantity",
	models.VProductOrderByUserName:  "user_name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
