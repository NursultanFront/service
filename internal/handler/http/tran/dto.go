package tran

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/models/types/money"
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/ardanlabs/service/internal/models/types/password"
	"github.com/ardanlabs/service/internal/models/types/quantity"
	"github.com/ardanlabs/service/internal/models/types/role"
)

// Product represents an individual product.
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

func toAppProduct(prd models.Product) Product {
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

// =============================================================================

// NewTran represents an example of cross domain transaction at the
// application layer.
type NewTran struct {
	Product NewProduct `json:"product"`
	User    NewUser    `json:"user"`
}

// Decode implements the decoder interface.
func (app *NewTran) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// =============================================================================

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Roles           []string `json:"roles"`
	Department      string   `json:"department"`
	Password        string   `json:"password"`
	PasswordConfirm string   `json:"passwordConfirm"`
}

func toBusNewUser(app NewUser) (models.NewUser, error) {
	var errors errs.FieldErrors

	roles, err := role.ParseMany(app.Roles)
	if err != nil {
		errors.Add("roles", err)
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		errors.Add("email", err)
	}

	nme, err := name.Parse(app.Name)
	if err != nil {
		errors.Add("name", err)
	}

	department, err := name.ParseNull(app.Department)
	if err != nil {
		errors.Add("department", err)
	}

	pass, err := password.ParseConfirm(app.Password, app.PasswordConfirm)
	if err != nil {
		errors.Add("password", err)
	}

	if len(errors) > 0 {
		return models.NewUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bus := models.NewUser{
		Name:       nme,
		Email:      *addr,
		Roles:      roles,
		Department: department,
		Password:   pass,
	}

	return bus, nil
}

// =============================================================================

// NewProduct is what we require from clients when adding a Product.
type NewProduct struct {
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Quantity int     `json:"quantity"`
}

func toBusNewProduct(app NewProduct) (models.NewProduct, error) {
	var errors errs.FieldErrors

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
		return models.NewProduct{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bus := models.NewProduct{
		Name:     name,
		Cost:     cost,
		Quantity: quantity,
	}

	return bus, nil
}
