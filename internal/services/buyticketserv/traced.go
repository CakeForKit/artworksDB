package buyticketserv

import (
	"context"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedBuyTicketsServ struct {
	buyTicketsServ BuyTicketsServ
	tracer         *tracing.Tracer
}

func NewTracedBuyTicketsServ(buyTicketsServ BuyTicketsServ, tracer *tracing.Tracer) BuyTicketsServ {
	return &tracedBuyTicketsServ{
		buyTicketsServ: buyTicketsServ,
		tracer:         tracer,
	}
}

func (b *tracedBuyTicketsServ) BuyTicket(
	ctx context.Context,
	eventID uuid.UUID,
	cntTickets int,
	customerName string,
	customerEmail string,
) (*models.TicketPurchaseTx, error) {
	var span trace.Span
	if b.tracer.IsEnabled() {
		ctx, span = b.tracer.StartSpan(ctx, "BuyTicket")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "buyticketserv"),
			attribute.String("event_id", eventID.String()),
			attribute.Int("ticket_count", cntTickets),
			attribute.String("customer_name", customerName),
			attribute.String("customer_email", customerEmail),
		)
	}

	return b.buyTicketsServ.BuyTicket(ctx, eventID, cntTickets, customerName, customerEmail)
}

func (b *tracedBuyTicketsServ) ConfirmBuyTicket(ctx context.Context, TxID uuid.UUID) error {
	var span trace.Span
	if b.tracer.IsEnabled() {
		ctx, span = b.tracer.StartSpan(ctx, "ConfirmBuyTicket")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "buyticketserv"),
			attribute.String("transaction_id", TxID.String()),
		)
	}

	return b.buyTicketsServ.ConfirmBuyTicket(ctx, TxID)
}

func (b *tracedBuyTicketsServ) CancelBuyTicket(ctx context.Context, TxID uuid.UUID) error {
	var span trace.Span
	if b.tracer.IsEnabled() {
		ctx, span = b.tracer.StartSpan(ctx, "CancelBuyTicket")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "buyticketserv"),
			attribute.String("transaction_id", TxID.String()),
		)
	}

	return b.buyTicketsServ.CancelBuyTicket(ctx, TxID)
}

func (b *tracedBuyTicketsServ) GetAllTicketPurchasesOfUser(ctx context.Context) ([]*models.TicketPurchase, error) {
	var span trace.Span
	if b.tracer.IsEnabled() {
		ctx, span = b.tracer.StartSpan(ctx, "GetAllTicketPurchasesOfUser")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "buyticketserv"),
		)
	}

	return b.buyTicketsServ.GetAllTicketPurchasesOfUser(ctx)
}

func (b *tracedBuyTicketsServ) GetBuyTicketTransactionDuration() time.Duration {
	var span trace.Span
	if b.tracer.IsEnabled() {
		_, span = b.tracer.StartSpan(context.Background(), "GetBuyTicketTransactionDuration")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "buyticketserv"),
		)
	}

	return b.buyTicketsServ.GetBuyTicketTransactionDuration()
}
