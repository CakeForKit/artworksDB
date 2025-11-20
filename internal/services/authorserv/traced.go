package authorserv

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedAuthorServ struct {
	authorServ AuthorServ
	tracer     *tracing.Tracer
}

func NewTracedAuthorServ(authorServ AuthorServ, tracer *tracing.Tracer) AuthorServ {
	return &tracedAuthorServ{
		authorServ: authorServ,
		tracer:     tracer,
	}
}

func (s *tracedAuthorServ) GetAll(ctx context.Context) ([]*models.Author, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "GetAllAuthors")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authorserv"),
		)
	}

	return s.authorServ.GetAll(ctx)
}

func (s *tracedAuthorServ) Add(ctx context.Context, author *models.Author) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "AddAuthor")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authorserv"),
			attribute.String("author_id", author.GetID().String()),
			attribute.String("author_name", author.GetName()),
		)
	}

	return s.authorServ.Add(ctx, author)
}

func (s *tracedAuthorServ) Update(ctx context.Context, idAuthor uuid.UUID, updateReq models.AuthorUpdateReq) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "UpdateAuthor")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authorserv"),
			attribute.String("author_id", idAuthor.String()),
			attribute.String("new_name", updateReq.Name),
		)
	}

	return s.authorServ.Update(ctx, idAuthor, updateReq)
}

func (s *tracedAuthorServ) Delete(ctx context.Context, idAuthor uuid.UUID) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "DeleteAuthor")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authorserv"),
			attribute.String("author_id", idAuthor.String()),
		)
	}

	return s.authorServ.Delete(ctx, idAuthor)
}
