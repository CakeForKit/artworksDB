package userservice

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedUserService struct {
	userService UserService
	tracer      *tracing.Tracer
}

func NewTracedUserService(userService UserService, tracer *tracing.Tracer) UserService {
	return &tracedUserService{
		userService: userService,
		tracer:      tracer,
	}
}

func (m *tracedUserService) ChangeSubscribeToMailing(ctx context.Context, subscr bool) error {
	var span trace.Span
	if m.tracer.IsEnabled() {
		ctx, span = m.tracer.StartSpan(ctx, "ChangeSubscribeToMailing")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "userservice"),
			attribute.Bool("subscription", subscr),
		)
	}

	return m.userService.ChangeSubscribeToMailing(ctx, subscr)
}

func (m *tracedUserService) GetSelf(ctx context.Context) (*models.User, error) {
	var span trace.Span
	if m.tracer.IsEnabled() {
		ctx, span = m.tracer.StartSpan(ctx, "GetSelf")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "userservice"),
		)
	}

	return m.userService.GetSelf(ctx)
}

func (m *tracedUserService) DeleteSelf(ctx context.Context) error {
	var span trace.Span
	if m.tracer.IsEnabled() {
		ctx, span = m.tracer.StartSpan(ctx, "DeleteSelf")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "userservice"),
		)
	}

	return m.userService.DeleteSelf(ctx)
}

func (m *tracedUserService) ChangePassword(ctx context.Context, newPassword string) error {
	var span trace.Span
	if m.tracer.IsEnabled() {
		ctx, span = m.tracer.StartSpan(ctx, "ChangePassword")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "userservice"),
		)
	}

	return m.userService.ChangePassword(ctx, newPassword)
}
