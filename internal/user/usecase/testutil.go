package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"net/mail"

	usermodels "github.com/ardanlabs/service/internal/user/models"

	"github.com/ardanlabs/service/internal/types/name"
	"github.com/ardanlabs/service/internal/types/password"
	"github.com/ardanlabs/service/internal/types/role"
	"github.com/google/uuid"
)

// TestNewUsers is a helper method for testing.
func TestNewUsers(n int, rle role.Role) []usermodels.NewUser {
	newUsrs := make([]usermodels.NewUser, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nu := usermodels.NewUser{
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
func TestSeedUsers(ctx context.Context, n int, role role.Role, api ExtBusiness) ([]usermodels.User, error) {
	newUsrs := TestNewUsers(n, role)

	usrs := make([]usermodels.User, len(newUsrs))
	for i, nu := range newUsrs {
		usr, err := api.Create(ctx, uuid.UUID{}, nu)
		if err != nil {
			return nil, fmt.Errorf("seeding user: idx: %d : %w", i, err)
		}

		usrs[i] = usr
	}

	return usrs, nil
}
