package audit_test

import (
	"context"
	"fmt"

	"github.com/ardanlabs/service/internal/app/sdk/apitest"
	"github.com/ardanlabs/service/internal/app/sdk/auth"
	"github.com/ardanlabs/service/internal/business/domain/auditbus"
	"github.com/ardanlabs/service/internal/business/domain/userbus"
	"github.com/ardanlabs/service/internal/business/sdk/dbtest"
	"github.com/ardanlabs/service/internal/business/types/domain"
	"github.com/ardanlabs/service/internal/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	audits, err := auditbus.TestSeedAudits(ctx, 2, usrs[0].ID, domain.User, "create", busDomain.Audit)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := apitest.User{
		User:   usrs[0],
		Audits: audits,
		Token:  apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{tu1},
	}

	return sd, nil
}
