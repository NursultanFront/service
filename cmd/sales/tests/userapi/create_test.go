package user_test

import (
	"net/http"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/handler/http/apitest"
	userapp "github.com/ardanlabs/service/internal/handler/http/user"
	"github.com/google/go-cmp/cmp"
)

func create200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &userapp.NewUser{
				Name:            "Bill Kennedy",
				Email:           "bill@ardanlabs.com",
				Roles:           []string{"ADMIN"},
				Department:      "ITO",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &userapp.User{},
			ExpResp: &userapp.User{
				Name:       "Bill Kennedy",
				Email:      "bill@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "ITO",
				Enabled:    true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*userapp.User)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*userapp.User)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &userapp.NewUser{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.InvalidArgument, "validate: [{\"field\":\"email\",\"error\":\"mail: no address\"},{\"field\":\"name\",\"error\":\"invalid name \\\"\\\"\"},{\"field\":\"password\",\"error\":\"invalid password \\\"\\\"\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-role",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &userapp.NewUser{
				Name:            "Bill Kennedy",
				Email:           "bill@ardanlabs.com",
				Roles:           []string{"SUPER"},
				Department:      "ITO",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Errorf(errs.InvalidArgument, "validate: [{\"field\":\"roles\",\"error\":\"invalid role \\\"SUPER\\\"\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-name",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &userapp.NewUser{
				Name:            "Bi",
				Email:           "bill@ardanlabs.com",
				Roles:           []string{"USER"},
				Department:      "ITO",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Errorf(errs.InvalidArgument, "validate: [{\"field\":\"name\",\"error\":\"invalid name \\\"Bi\\\"\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/users",
			Token:      "&nbsp;",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.Unauthenticated, "invalid or missing token"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badtoken",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token[:10],
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.Unauthenticated, "invalid or missing token"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token + "A",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.Unauthenticated, "invalid or missing token"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        "/v1/users",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.PermissionDenied, "you are not authorized for that action"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
