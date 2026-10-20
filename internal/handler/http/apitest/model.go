package apitest

import (
	auditmodels "github.com/ardanlabs/service/internal/audit/models"
	homemodels "github.com/ardanlabs/service/internal/home/models"
	productmodels "github.com/ardanlabs/service/internal/product/models"
	usermodels "github.com/ardanlabs/service/internal/user/models"
)

// User extends the dbtest user for api test support.
type User struct {
	usermodels.User
	Products []productmodels.Product
	Homes    []homemodels.Home
	Audits   []auditmodels.Audit
	Token    string
}

// SeedData represents users for api tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an api test.
type Table struct {
	Name       string
	URL        string
	Token      string
	Method     string
	StatusCode int
	Input      any
	GotResp    any
	ExpResp    any
	CmpFunc    func(got any, exp any) string
}
