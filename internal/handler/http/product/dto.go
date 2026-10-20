package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	productmodels "github.com/ardanlabs/service/internal/product/models"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/handler/http/mid"
	"github.com/ardanlabs/service/internal/types/money"
	"github.com/ardanlabs/service/internal/types/name"
	"github.com/ardanlabs/service/internal/types/quantity"
)

// Product represents information about an individual product.
type Product struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Product) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppProduct(prd productmodels.Product) Product {
	return Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProducts(prds []productmodels.Product) []Product {
	app := make([]Product, len(prds))
	for i, prd := range prds {
		app[i] = toAppProduct(prd)
	}

	return app
}

// =============================================================================

// NewProduct defines the data needed to add a new product.
type NewProduct struct {
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Quantity int     `json:"quantity"`
}

// Decode implements the decoder interface.
func (app *NewProduct) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewProduct(ctx context.Context, app NewProduct) (productmodels.NewProduct, error) {
	var errors errs.FieldErrors

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		errors.Add("userID", err)
	}

	name, err := name.Parse(app.Name)
	if err != nil {
		errors.Add("name", err)
	}

	cost, err := money.Parse(app.Cost)
	if err != nil {
		errors.Add("cost", err)
	}

	quantity, err := quantity.Parse(app.Quantity)
	if err != nil {
		errors.Add("quantity", err)
	}

	if len(errors) > 0 {
		return productmodels.NewProduct{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bus := productmodels.NewProduct{
		UserID:   userID,
		Name:     name,
		Cost:     cost,
		Quantity: quantity,
	}

	return bus, nil
}

// =============================================================================

// UpdateProduct defines the data needed to update a product.
type UpdateProduct struct {
	Name     *string  `json:"name"`
	Cost     *float64 `json:"cost"`
	Quantity *int     `json:"quantity"`
}

// Decode implements the decoder interface.
func (app *UpdateProduct) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusUpdateProduct(app UpdateProduct) (productmodels.UpdateProduct, error) {
	var errors errs.FieldErrors

	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			errors.Add("name", err)
		}
		nme = &nm
	}

	var cost *money.Money
	if app.Cost != nil {
		cst, err := money.Parse(*app.Cost)
		if err != nil {
			errors.Add("cost", err)
		}
		cost = &cst
	}

	var qnt *quantity.Quantity
	if app.Cost != nil {
		qn, err := quantity.Parse(*app.Quantity)
		if err != nil {
			errors.Add("quantity", err)
		}
		qnt = &qn
	}

	if len(errors) > 0 {
		return productmodels.UpdateProduct{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bus := productmodels.UpdateProduct{
		Name:     nme,
		Cost:     cost,
		Quantity: qnt,
	}

	return bus, nil
}
