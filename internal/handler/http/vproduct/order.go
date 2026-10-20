package vproduct

import (
	vproductbus "github.com/ardanlabs/service/internal/vproduct/usecase"
)

var orderByFields = map[string]string{
	"product_id": vproductbus.OrderByProductID,
	"user_id":    vproductbus.OrderByUserID,
	"name":       vproductbus.OrderByName,
	"cost":       vproductbus.OrderByCost,
	"quantity":   vproductbus.OrderByQuantity,
	"user_name":  vproductbus.OrderByUserName,
}
