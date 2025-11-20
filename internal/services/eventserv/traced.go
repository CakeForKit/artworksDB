package eventserv

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedEventService struct {
	eventService EventService
	tracer       *tracing.Tracer
}

func NewTracedEventService(eventService EventService, tracer *tracing.Tracer) EventService {
	return &tracedEventService{
		eventService: eventService,
		tracer:       tracer,
	}
}

func (e *tracedEventService) GetAll(ctx context.Context) ([]*models.Event, error) {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "GetAllEvents")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
		)
	}

	return e.eventService.GetAll(ctx)
}

func (e *tracedEventService) GetArtworksFromEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Artwork, error) {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "GetArtworksFromEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
			attribute.String("event_id", eventID.String()),
		)
	}

	return e.eventService.GetArtworksFromEvent(ctx, eventID)
}

func (e *tracedEventService) Add(ctx context.Context, eventReq *jsonreqresp.EventAdd) error {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "AddEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
			attribute.String("title", eventReq.Title),
			attribute.String("employee_id", eventReq.EmployeeID.String()),
			attribute.Int("ticket_count", eventReq.CntTickets),
			attribute.Bool("can_visit", eventReq.CanVisit),
			attribute.Int("artwork_count", len(eventReq.ArtworkIDs)),
		)
	}

	return e.eventService.Add(ctx, eventReq)
}

func (e *tracedEventService) Delete(ctx context.Context, eventID uuid.UUID) error {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "DeleteEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
			attribute.String("event_id", eventID.String()),
		)
	}

	return e.eventService.Delete(ctx, eventID)
}

func (e *tracedEventService) Update(
	ctx context.Context, eventID uuid.UUID, updateFields *jsonreqresp.EventUpdate) error {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "UpdateEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
			attribute.String("event_id", eventID.String()),
			attribute.String("title", updateFields.Title),
			attribute.String("address", updateFields.Address),
			attribute.Bool("can_visit", updateFields.CanVisit),
			attribute.Int("ticket_count", updateFields.CntTickets),
		)
	}

	return e.eventService.Update(ctx, eventID, updateFields)
}

func (e *tracedEventService) AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "AddArtworksToEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
			attribute.String("event_id", eventID.String()),
			attribute.Int("artwork_count", len(artworkIDs)),
		)
	}

	return e.eventService.AddArtworksToEvent(ctx, eventID, artworkIDs)
}

func (e *tracedEventService) DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "DeleteArtworkFromEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "eventserv"),
			attribute.String("event_id", eventID.String()),
			attribute.String("artwork_id", artworkID.String()),
		)
	}

	return e.eventService.DeleteArtworkFromEvent(ctx, eventID, artworkID)
}
