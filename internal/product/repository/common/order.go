package common

import (
	"fmt"

	productmodels "github.com/ardanlabs/service/internal/product/models"

	"github.com/ardanlabs/service/internal/sys/order"
)

// OrderByFields maps a productbus ordering key to its database column name.
// The mapping is shared because every engine uses the same column names.
var OrderByFields = map[string]string{
	productmodels.OrderByProductID: "product_id",
	productmodels.OrderByUserID:    "user_id",
	productmodels.OrderByName:      "name",
	productmodels.OrderByCost:      "cost",
	productmodels.OrderByQuantity:  "quantity",
}

// OrderByClause builds an ORDER BY clause for the given ordering. It
// returns an error if the field is not in OrderByFields.
func OrderByClause(orderBy order.By) (string, error) {
	by, exists := OrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
