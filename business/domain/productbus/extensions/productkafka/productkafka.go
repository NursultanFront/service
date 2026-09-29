// Package productkafka provides an extension for productbus that publishes
// product change events to Kafka.
//
// TODO(dev): this is a skeleton - implement Create/Update/Delete to actually
// marshal an event and write it via a *kafka.Writer (see
// business/sdk/kafkaclient for the constructor). Query/Count/QueryByID/
// QueryByUserID are pure passthrough (no event to publish for reads) and
// don't need changes.
package productkafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ardanlabs/service/business/domain/productbus"
	"github.com/ardanlabs/service/business/sdk/order"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Extension provides a wrapper that publishes product change events to
// Kafka around the productbus.
type Extension struct {
	log    *logger.Logger
	bus    productbus.ExtBusiness
	writer *kafka.Writer
}

type productEvent struct {
	ProductID uuid.UUID `json:"product_id"`
	EventType string    `json:"event_type"`
	Message   string    `json:"message"`
}

// NewExtension constructs a new extension that wraps the productbus with a
// Kafka producer. writer is expected to already be configured for the
// target topic (see kafkaclient.NewWriter).
//
// NOTE(AI): log field/param added by Claude by direct request, one-off
// exception per AGENTS.md - so publish failures can be logged without
// failing the (already-committed) Create/Update/Delete call.
func NewExtension(log *logger.Logger, writer *kafka.Writer) productbus.Extension {
	return func(bus productbus.ExtBusiness) productbus.ExtBusiness {
		return &Extension{
			log:    log,
			bus:    bus,
			writer: writer,
		}
	}
}

// NewWithTx does not publish events - a transactional write is not yet
// committed and may be rolled back.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (productbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Create publishes a product-created event after a successful create.
//
// NOTE(AI): разработчик написал структуру события и вызовы Marshal/
// WriteMessages сам; Claude по прямой просьбе поправил недостающий импорт
// encoding/json, значения EventType/Message и имя переменной (совпадало с
// именем типа productEvent) - разовое исключение из AGENTS.md.
func (ext *Extension) Create(ctx context.Context, np productbus.NewProduct) (productbus.Product, error) {
	prd, err := ext.bus.Create(ctx, np)
	if err != nil {
		return productbus.Product{}, err
	}

	event := productEvent{
		ProductID: prd.ID,
		EventType: "product.created",
		Message:   fmt.Sprintf("product %s created", prd.Name.String()),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return productbus.Product{}, err
	}

	if err := ext.writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		ext.log.Error(ctx, "productkafka: write", "ERROR", err)
	}

	return prd, nil
}

// Update publishes a product-updated event after a successful update.
//
// TODO(dev): same as Create, but for an update event.
func (ext *Extension) Update(ctx context.Context, prd productbus.Product, up productbus.UpdateProduct) (productbus.Product, error) {
	updated, err := ext.bus.Update(ctx, prd, up)
	if err != nil {
		return productbus.Product{}, err
	}

	event := productEvent{
		ProductID: updated.ID,
		EventType: "product.updated",
		Message:   fmt.Sprintf("product %s updated", updated.Name.String()),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return productbus.Product{}, err
	}

	if err := ext.writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		ext.log.Error(ctx, "productkafka: write", "ERROR", err)
	}

	return updated, nil
}

// Delete publishes a product-deleted event after a successful delete.
//
// TODO(dev): same as Create, but for a delete event.
func (ext *Extension) Delete(ctx context.Context, prd productbus.Product) error {
	if err := ext.bus.Delete(ctx, prd); err != nil {
		return err
	}

	event := productEvent{
		ProductID: prd.ID,
		EventType: "product.deleted",
		Message:   fmt.Sprintf("product %s deleted", prd.Name.String()),
	}

	data, err := json.Marshal(event)
	
	if err != nil {
		return err
	}

	if err := ext.writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		ext.log.Error(ctx, "productkafka: write", "ERROR", err)
	}

	return nil
}

// Query is a pure passthrough - no event to publish for a read.
func (ext *Extension) Query(ctx context.Context, filter productbus.QueryFilter, orderBy order.By, page page.Page) ([]productbus.Product, error) {
	return ext.bus.Query(ctx, filter, orderBy, page)
}

// Count is a pure passthrough - no event to publish for a read.
func (ext *Extension) Count(ctx context.Context, filter productbus.QueryFilter) (int, error) {
	return ext.bus.Count(ctx, filter)
}

// QueryByID is a pure passthrough - no event to publish for a read.
func (ext *Extension) QueryByID(ctx context.Context, productID uuid.UUID) (productbus.Product, error) {
	return ext.bus.QueryByID(ctx, productID)
}

// QueryByUserID is a pure passthrough - no event to publish for a read.
func (ext *Extension) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]productbus.Product, error) {
	return ext.bus.QueryByUserID(ctx, userID)
}
