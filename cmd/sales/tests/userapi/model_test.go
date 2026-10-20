package user_test

import (
	"time"

	usermodels "github.com/ardanlabs/service/internal/user/models"

	userapp "github.com/ardanlabs/service/internal/handler/http/user"
	"github.com/ardanlabs/service/internal/types/role"
)

func toAppUser(bus usermodels.User) userapp.User {
	return userapp.User{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Email:       bus.Email.Address,
		Roles:       role.ParseToString(bus.Roles),
		Department:  bus.Department.String(),
		Enabled:     bus.Enabled,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []usermodels.User) []userapp.User {
	items := make([]userapp.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(bus usermodels.User) *userapp.User {
	appUsr := toAppUser(bus)
	return &appUsr
}
