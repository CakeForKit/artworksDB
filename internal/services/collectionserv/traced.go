package collectionserv

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedCollectionServ struct {
	collectionServ CollectionServ
	tracer         *tracing.Tracer
}

func NewTracedCollectionServ(collectionServ CollectionServ, tracer *tracing.Tracer) CollectionServ {
	return &tracedCollectionServ{
		collectionServ: collectionServ,
		tracer:         tracer,
	}
}

func (s *tracedCollectionServ) GetAll(ctx context.Context) ([]*models.Collection, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetAllCollections")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "collectionserv"),
		)
	}

	return s.collectionServ.GetAll(ctx)
}

func (s *tracedCollectionServ) Add(ctx context.Context, col *models.Collection) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "AddCollection")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "collectionserv"),
			attribute.String("collection_id", col.GetID().String()),
			attribute.String("collection_name", col.GetTitle()),
		)
	}

	return s.collectionServ.Add(ctx, col)
}

func (s *tracedCollectionServ) Update(
	ctx context.Context, idCol uuid.UUID, updateReq models.CollectionUpdateReq) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "UpdateCollection")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "collectionserv"),
			attribute.String("collection_id", idCol.String()),
			attribute.String("new_name", updateReq.Title),
		)
	}

	return s.collectionServ.Update(ctx, idCol, updateReq)
}

func (s *tracedCollectionServ) Delete(ctx context.Context, idCol uuid.UUID) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "DeleteCollection")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "collectionserv"),
			attribute.String("collection_id", idCol.String()),
		)
	}

	return s.collectionServ.Delete(ctx, idCol)
}
