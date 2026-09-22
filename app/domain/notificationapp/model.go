package notificationapp

import (
	"encoding/json"
	"time"

	"github.com/ardanlabs/service/business/domain/notificationbus"
)

// Notification represents information about an individual notification.
type Notification struct {
	ID        string `json:"id"`
	ProductID string `json:"productID"`
	EventType string `json:"eventType"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

// Notifications represents a list of notifications for HTTP responses.
type Notifications []Notification

// Encode implements the encoder interface.
func (app Notifications) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppNotification(n notificationbus.Notification) Notification {
	return Notification{
		ID:        n.ID.String(),
		ProductID: n.ProductID.String(),
		EventType: n.EventType,
		Message:   n.Message,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
}

func toAppNotifications(ns []notificationbus.Notification) Notifications {
	items := make(Notifications, len(ns))
	for i, n := range ns {
		items[i] = toAppNotification(n)
	}

	return items
}
