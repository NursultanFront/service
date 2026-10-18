package user

import (
	"context"
	"fmt"
	"math/rand"
	"net/mail"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/ardanlabs/service/internal/models/types/password"
	"github.com/ardanlabs/service/internal/models/types/role"
	"github.com/google/uuid"
)

// TestNewUsers is a helper method for testing.
func TestNewUsers(n int, rle role.Role) []models.NewUser {
	newUsrs := make([]models.NewUser, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nu := models.NewUser{
			Name:       name.MustParse(fmt.Sprintf("Name%d", idx)),
			Email:      mail.Address{Address: fmt.Sprintf("Email%d@gmail.com", idx)},
			Roles:      []role.Role{rle},
			Department: name.MustParseNull(fmt.Sprintf("Department%d", idx)),
			Password:   password.MustParse(fmt.Sprintf("Password%d", idx)),
		}

		newUsrs[i] = nu
	}

	return newUsrs
}

// TestSeedUsers is a helper method for testing.
func TestSeedUsers(ctx context.Context, n int, role role.Role, api ExtBusiness) ([]models.User, error) {
	newUsrs := TestNewUsers(n, role)

	usrs := make([]models.User, len(newUsrs))
	for i, nu := range newUsrs {
		usr, err := api.Create(ctx, uuid.UUID{}, nu)
		if err != nil {
			return nil, fmt.Errorf("seeding user: idx: %d : %w", i, err)
		}

		usrs[i] = usr
	}

	return usrs, nil
}
