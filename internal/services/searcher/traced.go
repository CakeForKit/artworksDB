package searcher

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedSearcher struct {
	searcher Searcher
	tracer   *tracing.Tracer
}

func NewTracedSearcher(searcher Searcher, tracer *tracing.Tracer) Searcher {
	return &tracedSearcher{
		searcher: searcher,
		tracer:   tracer,
	}
}

func (s *tracedSearcher) GetAllArtworks(ctx context.Context,
	filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetAllArtworks")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "searcher"),
		)
	}

	return s.searcher.GetAllArtworks(ctx, filterOps, sortOps)
}

func (s *tracedSearcher) GetAllEvents(
	ctx context.Context, filterOps *jsonreqresp.EventFilter,
) ([]*models.Event, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetAllEvents")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "searcher"),
		)
	}

	return s.searcher.GetAllEvents(ctx, filterOps)
}

func (s *tracedSearcher) GetEvent(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "searcher"),
		)
	}

	return s.searcher.GetEvent(ctx, eventID)
}

func (s *tracedSearcher) GetArtworksFromEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Artwork, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetArtworksFromEvent")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "searcher"),
		)
	}

	return s.searcher.GetArtworksFromEvent(ctx, eventID)
}

func (s *tracedSearcher) GetCollectionsStat(ctx context.Context, eventID uuid.UUID) ([]*models.StatCollections, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetCollectionsStat")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "searcher"),
		)
	}

	return s.searcher.GetCollectionsStat(ctx, eventID)
}
