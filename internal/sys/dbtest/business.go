package dbtest

import (
	"time"

	auditdb "github.com/ardanlabs/service/internal/audit/repository"
	auditbus "github.com/ardanlabs/service/internal/audit/usecase"
	"github.com/ardanlabs/service/internal/audit/usecase/extensions/auditotel"
	homedb "github.com/ardanlabs/service/internal/home/repository"
	homebus "github.com/ardanlabs/service/internal/home/usecase"
	"github.com/ardanlabs/service/internal/home/usecase/extensions/homeotel"
	productpg "github.com/ardanlabs/service/internal/product/repository/postgres"
	productbus "github.com/ardanlabs/service/internal/product/usecase"
	"github.com/ardanlabs/service/internal/product/usecase/extensions/productotel"
	"github.com/ardanlabs/service/internal/sys/delegate"
	userdb "github.com/ardanlabs/service/internal/user/repository"
	usercache "github.com/ardanlabs/service/internal/user/repository/cache"
	userbus "github.com/ardanlabs/service/internal/user/usecase"
	"github.com/ardanlabs/service/internal/user/usecase/extensions/useraudit"
	"github.com/ardanlabs/service/internal/user/usecase/extensions/userotel"
	vproductdb "github.com/ardanlabs/service/internal/vproduct/repository"
	vproductbus "github.com/ardanlabs/service/internal/vproduct/usecase"
	"github.com/ardanlabs/service/internal/vproduct/usecase/extensions/vproductotel"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	Audit    auditbus.ExtBusiness
	Home     homebus.ExtBusiness
	Product  productbus.ExtBusiness
	User     userbus.ExtBusiness
	VProduct vproductbus.ExtBusiness
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)

	auditOtelExt := auditotel.NewExtension()
	auditStorage := auditdb.NewStore(log, db)
	auditBus := auditbus.NewBusiness(log, auditStorage, auditOtelExt)

	userOtelExt := userotel.NewExtension()
	userAuditExt := useraudit.NewExtension(auditBus)
	userStorage := usercache.NewStore(log, userdb.NewStore(log, db), nil, time.Hour)
	userBus := userbus.NewBusiness(log, delegate, userStorage, userOtelExt, userAuditExt)

	productOtelExt := productotel.NewExtension()
	productStorage := productpg.NewStore(log, db)
	productBus := productbus.NewBusiness(log, userBus, delegate, productStorage, productOtelExt)

	homeOtelExt := homeotel.NewExtension()
	homeStorage := homedb.NewStore(log, db)
	homeBus := homebus.NewBusiness(log, userBus, delegate, homeStorage, homeOtelExt)

	vproductOtelExt := vproductotel.NewExtension()
	vproductStorage := vproductdb.NewStore(log, db)
	vproductBus := vproductbus.NewBusiness(vproductStorage, vproductOtelExt)

	return BusDomain{
		Delegate: delegate,
		Audit:    auditBus,
		Home:     homeBus,
		Product:  productBus,
		User:     userBus,
		VProduct: vproductBus,
	}
}
