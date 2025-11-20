package authadmin

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedAuthAdmin struct {
	authAdmin AdminAuth
	tracer    *tracing.Tracer
}

func NewTracedAuthAdmin(authAdmin AdminAuth, tracer *tracing.Tracer) AdminAuth {
	return &tracedAuthAdmin{
		authAdmin: authAdmin,
		tracer:    tracer,
	}
}

func (s *tracedAuthAdmin) LoginAdmin(ctx context.Context, lur LoginAdminRequest) (string, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "LoginAdmin")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "auth"),
			attribute.String("login", lur.Login),
			attribute.String("role", "admin"),
		)
	}

	return s.authAdmin.LoginAdmin(ctx, lur)
}

func (s *tracedAuthAdmin) RegisterAdmin(ctx context.Context, rur RegisterAdminRequest) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "RegisterAdmin")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "auth"),
			attribute.String("adminname", rur.Adminname),
			attribute.String("login", rur.Login),
			attribute.String("role", "admin"),
		)
	}

	return s.authAdmin.RegisterAdmin(ctx, rur)
}

func (s *tracedAuthAdmin) VerifyByToken(tokenStr string) (*token.Payload, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		_, span = s.tracer.StartSpan(context.Background(), "VerifyByToken")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "auth"),
			attribute.String("role", "admin"),
		)
	}

	return s.authAdmin.VerifyByToken(tokenStr)
}
