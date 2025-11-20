package authuser

import (
	"context"

	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedAuthUser struct {
	authUser AuthUser
	tracer   *tracing.Tracer
}

func NewTracedAuthUser(authUser AuthUser, tracer *tracing.Tracer) AuthUser {
	return &tracedAuthUser{
		authUser: authUser,
		tracer:   tracer,
	}
}

func (s *tracedAuthUser) RegisterUser(ctx context.Context, rur authmodels.RegisterUserRequest) error {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "RegisterUser")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authuser"),
			attribute.String("username", rur.Username),
			attribute.String("login", rur.Login),
			attribute.String("email", rur.Email),
			attribute.Bool("subscribe_email", rur.SubscribeEmail),
			attribute.String("role", "user"),
		)
	}

	return s.authUser.RegisterUser(ctx, rur)
}

func (s *tracedAuthUser) LoginUser(ctx context.Context, lur authmodels.LoginUserRequest) (uuid.UUID, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "LoginUser")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authuser"),
			attribute.String("login", lur.Login),
			attribute.String("role", "user"),
		)
	}

	return s.authUser.LoginUser(ctx, lur)
}

func (s *tracedAuthUser) OTP(ctx context.Context, sessionID uuid.UUID, otpCode string) (string, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		ctx, span = s.tracer.StartSpan(ctx, "OTP")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authuser"),
			attribute.String("session_id", sessionID.String()),
			attribute.String("role", "user"),
		)
	}

	return s.authUser.OTP(ctx, sessionID, otpCode)
}

func (s *tracedAuthUser) VerifyByToken(tokenStr string) (*token.Payload, error) {
	var span trace.Span
	if s.tracer.IsEnabled() {
		_, span = s.tracer.StartSpan(context.Background(), "VerifyByToken")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "authuser"),
			attribute.String("role", "user"),
		)
	}

	return s.authUser.VerifyByToken(tokenStr)
}
