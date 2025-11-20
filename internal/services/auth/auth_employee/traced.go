package authemployee

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedAuthEmployee struct {
	authEmployee AuthEmployee
	tracer       *tracing.Tracer
}

func NewTracedAuthEmployee(authEmployee AuthEmployee, tracer *tracing.Tracer) AuthEmployee {
	return &tracedAuthEmployee{
		authEmployee: authEmployee,
		tracer:       tracer,
	}
}

func (s *tracedAuthEmployee) LoginEmployee(ctx context.Context, ler LoginEmployeeRequest) (string, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "LoginEmployee")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authemployee"),
			attribute.String("login", ler.Login),
			attribute.String("role", "employee"),
		)
	}

	return s.authEmployee.LoginEmployee(ctx, ler)
}

func (s *tracedAuthEmployee) RegisterEmployee(
	ctx context.Context, rer RegisterEmployeeRequest, adminID uuid.UUID,
) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "RegisterEmployee")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authemployee"),
			attribute.String("username", rer.Username),
			attribute.String("login", rer.Login),
			attribute.String("admin_id", adminID.String()),
			attribute.String("role", "employee"),
		)
	}

	return s.authEmployee.RegisterEmployee(ctx, rer, adminID)
}

func (s *tracedAuthEmployee) VerifyByToken(tokenStr string) (*token.Payload, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		_, span = s.tracer.StartSpan(context.Background(), "VerifyByToken")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authemployee"),
			attribute.String("role", "employee"),
		)
	}

	return s.authEmployee.VerifyByToken(tokenStr)
}
