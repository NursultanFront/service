package repository

import (
	"fmt"

	vproductmodels "github.com/ardanlabs/service/internal/vproduct/models"

	"github.com/ardanlabs/service/internal/sys/order"
)

var orderByFields = map[string]string{
	vproductmodels.OrderByProductID: "product_id",
	vproductmodels.OrderByUserID:    "user_id",
	vproductmodels.OrderByName:      "name",
	vproductmodels.OrderByCost:      "cost",
	vproductmodels.OrderByQuantity:  "quantity",
	vproductmodels.OrderByUserName:  "user_name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
