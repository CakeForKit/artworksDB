package authz

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedAuthZ struct {
	authZ  AuthZ
	tracer *tracing.Tracer
}

func NewTracedAuthZ(authZ AuthZ, tracer *tracing.Tracer) AuthZ {
	return &tracedAuthZ{
		authZ:  authZ,
		tracer: tracer,
	}
}

func (a *tracedAuthZ) Authorize(ctx context.Context, payload token.Payload) context.Context {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "Authorize")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authz"),
			attribute.String("person_id", payload.PersonID.String()),
			attribute.String("role", payload.Role),
		)
	}

	return a.authZ.Authorize(ctx, payload)
}

func (a *tracedAuthZ) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "UserIDFromContext")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authz"),
			attribute.String("role", "user"),
		)
	}

	return a.authZ.UserIDFromContext(ctx)
}

func (a *tracedAuthZ) EmployeeIDFromContext(ctx context.Context) (uuid.UUID, error) {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "EmployeeIDFromContext")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authz"),
			attribute.String("role", "employee"),
		)
	}

	return a.authZ.EmployeeIDFromContext(ctx)
}

func (a *tracedAuthZ) AdminIDFromContext(ctx context.Context) (uuid.UUID, error) {
	var span trace.Span
	if a.tracer.IsEnabled() {
		ctx, span = a.tracer.StartSpan(ctx, "AdminIDFromContext")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authz"),
			attribute.String("role", "admin"),
		)
	}

	return a.authZ.AdminIDFromContext(ctx)
}
