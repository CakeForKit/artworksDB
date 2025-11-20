package artworkserv

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedArtworkService struct {
	artworkService ArtworkService
	tracer         *tracing.Tracer
}

func NewTracedArtworkService(artworkService ArtworkService, tracer *tracing.Tracer) ArtworkService {
	return &tracedArtworkService{
		artworkService: artworkService,
		tracer:         tracer,
	}
}

func (a *tracedArtworkService) GetAll(ctx context.Context) ([]*models.Artwork, error) {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "GetAllArtworks")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "artworkserv"),
		)
	}

	return a.artworkService.GetAll(ctx)
}

func (a *tracedArtworkService) Add(ctx context.Context, artworkReq jsonreqresp.AddArtworkRequest) error {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "AddArtwork")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "artworkserv"),
			attribute.String("title", artworkReq.Title),
			attribute.String("author_id", artworkReq.AuthorID),
			attribute.String("collection_id", artworkReq.CollectionID),
			attribute.String("technic", artworkReq.Technic),
			attribute.String("material", artworkReq.Material),
			attribute.Int("creation_year", artworkReq.CreationYear),
		)
	}

	return a.artworkService.Add(ctx, artworkReq)
}

func (a *tracedArtworkService) Delete(ctx context.Context, idArt uuid.UUID) error {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "DeleteArtwork")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "artworkserv"),
			attribute.String("artwork_id", idArt.String()),
		)
	}

	return a.artworkService.Delete(ctx, idArt)
}

func (a *tracedArtworkService) Update(
	ctx context.Context, idArt uuid.UUID, updateFields jsonreqresp.ArtworkUpdate) error {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "UpdateArtwork")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "artworkserv"),
			attribute.String("artwork_id", idArt.String()),
			attribute.String("author_id", updateFields.AuthorID),
			attribute.String("collection_id", updateFields.CollectionID),
			attribute.String("title", updateFields.Title),
			attribute.String("technic", updateFields.Technic),
			attribute.String("material", updateFields.Material),
			attribute.Int("creation_year", updateFields.CreationYear),
		)
	}

	return a.artworkService.Update(ctx, idArt, updateFields)
}
