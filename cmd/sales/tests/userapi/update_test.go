package user_test

import (
	"fmt"
	"net/http"
	"time"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/handler/http/apitest"
	userapp "github.com/ardanlabs/service/internal/handler/http/user"
	"github.com/google/go-cmp/cmp"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &userapp.UpdateUser{
				Name:            new("Jack Kennedy"),
				Email:           new("jack@ardanlabs.com"),
				Department:      new("ITO"),
				Password:        new("123"),
				PasswordConfirm: new("123"),
			},
			GotResp: &userapp.User{},
			ExpResp: &userapp.User{
				ID:          sd.Users[0].ID.String(),
				Name:        "Jack Kennedy",
				Email:       "jack@ardanlabs.com",
				Roles:       []string{"USER"},
				Department:  "ITO",
				Enabled:     true,
				DateCreated: sd.Users[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].DateUpdated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*userapp.User)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*userapp.User)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "role",
			URL:        fmt.Sprintf("/v1/users/role/%s", sd.Admins[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &userapp.UpdateUserRole{
				Roles: []string{"USER"},
			},
			GotResp: &userapp.User{},
			ExpResp: &userapp.User{
				ID:          sd.Admins[0].ID.String(),
				Name:        sd.Admins[0].Name.String(),
				Email:       sd.Admins[0].Email.Address,
				Roles:       []string{"USER"},
				Department:  sd.Admins[0].Department.String(),
				Enabled:     true,
				DateCreated: sd.Admins[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Admins[0].DateUpdated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*userapp.User)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*userapp.User)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bad-input",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &userapp.UpdateUser{
				Email:           new("bill@"),
				PasswordConfirm: new("jack"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Errorf(errs.InvalidArgument, "validate: [{\"field\":\"email\",\"error\":\"mail: missing '@' or angle-addr\"},{\"field\":\"password\",\"error\":\"passwords do not match\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-role",
			URL:        fmt.Sprintf("/v1/users/role/%s", sd.Admins[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &userapp.UpdateUserRole{
				Roles: []string{"BAD ROLE"},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Errorf(errs.InvalidArgument, "validate: [{\"field\":\"roles\",\"error\":\"invalid role \\\"BAD ROLE\\\"\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      "&nbsp;",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.Unauthenticated, "invalid or missing token"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token + "A",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Errorf(errs.Unauthenticated, "invalid or missing token"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Admins[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &userapp.UpdateUser{
				Name:            new("Bill Kennedy"),
				Email:           new("bill@ardanlabs.com"),
				Department:      new("ITO"),
				Password:        new("123"),
				PasswordConfirm: new("123"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Errorf(errs.PermissionDenied, "you are not authorized for that action"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "roleadminonly",
			URL:        fmt.Sprintf("/v1/users/role/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &userapp.UpdateUserRole{
				Roles: []string{"ADMIN"},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Errorf(errs.PermissionDenied, "you are not authorized for that action"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
