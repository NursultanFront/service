package product_test

import (
	"time"

	productapp "github.com/ardanlabs/service/internal/handler/http/product"
	"github.com/ardanlabs/service/internal/models"
)

func toAppProduct(prd models.Product) productapp.Product {
	return productapp.Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProductPtr(prd models.Product) *productapp.Product {
	appPrd := toAppProduct(prd)
	return &appPrd
}

func toAppProducts(prds []models.Product) []productapp.Product {
	items := make([]productapp.Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
