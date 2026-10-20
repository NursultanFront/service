package product_test

import (
	"time"

	productmodels "github.com/ardanlabs/service/internal/product/models"

	productapp "github.com/ardanlabs/service/internal/handler/http/product"
)

func toAppProduct(prd productmodels.Product) productapp.Product {
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

func toAppProductPtr(prd productmodels.Product) *productapp.Product {
	appPrd := toAppProduct(prd)
	return &appPrd
}

func toAppProducts(prds []productmodels.Product) []productapp.Product {
	items := make([]productapp.Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
