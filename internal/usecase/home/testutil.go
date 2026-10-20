package home

import (
	"context"
	"fmt"
	"math/rand"

	homemodels "github.com/ardanlabs/service/internal/home/models"

	"github.com/ardanlabs/service/internal/types/home"
	"github.com/google/uuid"
)

// TestGenerateNewHomes is a helper method for testing.
func TestGenerateNewHomes(n int, userID uuid.UUID) []homemodels.NewHome {
	newHmes := make([]homemodels.NewHome, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		t := home.Single
		if v := (idx + i) % 2; v == 0 {
			t = home.Condo
		}

		nh := homemodels.NewHome{
			Type: t,
			Address: homemodels.Address{
				Address1: fmt.Sprintf("Address%d", idx),
				Address2: fmt.Sprintf("Address%d", idx),
				ZipCode:  fmt.Sprintf("%05d", idx),
				City:     fmt.Sprintf("City%d", idx),
				State:    fmt.Sprintf("State%d", idx),
				Country:  fmt.Sprintf("Country%d", idx),
			},
			UserID: userID,
		}

		newHmes[i] = nh
	}

	return newHmes
}

// TestGenerateSeedHomes is a helper method for testing.
func TestGenerateSeedHomes(ctx context.Context, n int, api ExtBusiness, userID uuid.UUID) ([]homemodels.Home, error) {
	newHmes := TestGenerateNewHomes(n, userID)

	hmes := make([]homemodels.Home, len(newHmes))
	for i, nh := range newHmes {
		hme, err := api.Create(ctx, nh)
		if err != nil {
			return nil, fmt.Errorf("seeding home: idx: %d : %w", i, err)
		}

		hmes[i] = hme
	}

	return hmes, nil
}

// ParseAddress is a helper function to create an address value.
func ParseAddress(address1 string, address2 string, zipCode string, city string, state string, country string) homemodels.Address {
	return homemodels.Address{
		Address1: address1,
		Address2: address2,
		ZipCode:  zipCode,
		City:     city,
		State:    state,
		Country:  country,
	}
}
