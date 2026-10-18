package home

import (
	homebus "github.com/ardanlabs/service/internal/usecase/home"
)

var orderByFields = map[string]string{
	"home_id": homebus.OrderByID,
	"type":    homebus.OrderByType,
	"user_id": homebus.OrderByUserID,
}
