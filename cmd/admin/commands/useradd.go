package commands

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	usermodels "github.com/ardanlabs/service/internal/user/models"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/types/name"
	"github.com/ardanlabs/service/internal/types/password"
	"github.com/ardanlabs/service/internal/types/role"
	userdb "github.com/ardanlabs/service/internal/user/repository"
	userbus "github.com/ardanlabs/service/internal/user/usecase"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/google/uuid"
)

// UserAdd adds new users into the database.
func UserAdd(log *logger.Logger, cfg sqldb.Config, nme string, email string, pass string) error {
	if nme == "" || email == "" || pass == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewBusiness(log, nil, userdb.NewStore(log, db))

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := usermodels.NewUser{
		Name:     name.MustParse(nme),
		Email:    *addr,
		Password: password.MustParse(pass),
		Roles:    []role.Role{role.Admin, role.User},
	}

	usr, err := userBus.Create(ctx, uuid.UUID{}, nu)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
