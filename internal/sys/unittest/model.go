package unittest

import (
	"context"

	auditmodels "github.com/ardanlabs/service/internal/audit/models"
	homemodels "github.com/ardanlabs/service/internal/home/models"
	productmodels "github.com/ardanlabs/service/internal/product/models"
	usermodels "github.com/ardanlabs/service/internal/user/models"
)

// User represents an app user specified for the test.
type User struct {
	usermodels.User
	Products []productmodels.Product
	Homes    []homemodels.Home
	Audits   []auditmodels.Audit
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an unit test.
type Table struct {
	Name    string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}
