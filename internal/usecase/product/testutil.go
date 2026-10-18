package product

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/models/types/money"
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/ardanlabs/service/internal/models/types/quantity"
	"github.com/google/uuid"
)

// TestGenerateNewProducts is a helper method for testing.
func TestGenerateNewProducts(n int, userID uuid.UUID) []models.NewProduct {
	newPrds := make([]models.NewProduct, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		np := models.NewProduct{
			Name:     name.MustParse(fmt.Sprintf("Name%d", idx)),
			Cost:     money.MustParse(float64(rand.Intn(500) + 1)),
			Quantity: quantity.MustParse(rand.Intn(50) + 1),
			UserID:   userID,
		}

		newPrds[i] = np
	}

	return newPrds
}

// TestGenerateSeedProducts is a helper method for testing.
func TestGenerateSeedProducts(ctx context.Context, n int, api ExtBusiness, userID uuid.UUID) ([]models.Product, error) {
	newPrds := TestGenerateNewProducts(n, userID)

	prds := make([]models.Product, len(newPrds))
	for i, np := range newPrds {
		prd, err := api.Create(ctx, np)
		if err != nil {
			return nil, fmt.Errorf("seeding product: idx: %d : %w", i, err)
		}

		prds[i] = prd
	}

	return prds, nil
}
