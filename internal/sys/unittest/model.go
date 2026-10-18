package unittest

import (
	"context"

	"github.com/ardanlabs/service/internal/models"
)

// User represents an app user specified for the test.
type User struct {
	models.User
	Products []models.Product
	Homes    []models.Home
	Audits   []models.Audit
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
