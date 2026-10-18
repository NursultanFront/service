package vproduct_test

import (
	"time"

	vproductapp "github.com/ardanlabs/service/internal/handler/http/vproduct"
	"github.com/ardanlabs/service/internal/models"
)

func toAppVProduct(usr models.User, prd models.Product) vproductapp.Product {
	return vproductapp.Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
		UserName:    usr.Name.String(),
	}
}

func toAppVProducts(usr models.User, prds []models.Product) []vproductapp.Product {
	items := make([]vproductapp.Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppVProduct(usr, prd)
	}

	return items
}
