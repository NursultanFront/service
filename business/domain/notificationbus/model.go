package notificationbus

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents an individual notification, typically created
// from a Kafka event (e.g. a product change).
//
// TODO(dev): shape this to whatever the Kafka event you publish from
// productkafka actually looks like.
type Notification struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	EventType string
	Message   string
	CreatedAt time.Time
}

// NewNotification contains the information needed to create a new
// Notification.
type NewNotification struct {
	ProductID uuid.UUID
	EventType string
	Message   string
}
