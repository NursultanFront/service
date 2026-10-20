package home_test

import (
	"time"

	homemodels "github.com/ardanlabs/service/internal/home/models"

	homeapp "github.com/ardanlabs/service/internal/handler/http/home"
)

func toAppHome(hme homemodels.Home) homeapp.Home {
	return homeapp.Home{
		ID:     hme.ID.String(),
		UserID: hme.UserID.String(),
		Type:   hme.Type.String(),
		Address: homeapp.Address{
			Address1: hme.Address.Address1,
			Address2: hme.Address.Address2,
			ZipCode:  hme.Address.ZipCode,
			City:     hme.Address.City,
			State:    hme.Address.State,
			Country:  hme.Address.Country,
		},
		DateCreated: hme.DateCreated.Format(time.RFC3339),
		DateUpdated: hme.DateUpdated.Format(time.RFC3339),
	}
}

func toAppHomes(homes []homemodels.Home) []homeapp.Home {
	items := make([]homeapp.Home, len(homes))
	for i, hme := range homes {
		items[i] = toAppHome(hme)
	}

	return items
}

func toAppHomePtr(hme homemodels.Home) *homeapp.Home {
	appHme := toAppHome(hme)
	return &appHme
}
