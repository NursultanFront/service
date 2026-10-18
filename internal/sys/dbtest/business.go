package dbtest

import (
	"time"

	auditdb "github.com/ardanlabs/service/internal/repository/audit"
	homedb "github.com/ardanlabs/service/internal/repository/home"
	productpg "github.com/ardanlabs/service/internal/repository/product/postgres"
	userdb "github.com/ardanlabs/service/internal/repository/user"
	usercache "github.com/ardanlabs/service/internal/repository/user/cache"
	vproductdb "github.com/ardanlabs/service/internal/repository/vproduct"
	"github.com/ardanlabs/service/internal/sys/delegate"
	auditbus "github.com/ardanlabs/service/internal/usecase/audit"
	"github.com/ardanlabs/service/internal/usecase/audit/extensions/auditotel"
	homebus "github.com/ardanlabs/service/internal/usecase/home"
	"github.com/ardanlabs/service/internal/usecase/home/extensions/homeotel"
	productbus "github.com/ardanlabs/service/internal/usecase/product"
	"github.com/ardanlabs/service/internal/usecase/product/extensions/productotel"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
	"github.com/ardanlabs/service/internal/usecase/user/extensions/useraudit"
	"github.com/ardanlabs/service/internal/usecase/user/extensions/userotel"
	vproductbus "github.com/ardanlabs/service/internal/usecase/vproduct"
	"github.com/ardanlabs/service/internal/usecase/vproduct/extensions/vproductotel"
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
