package authz_test

import (
	"context"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/services/auth/authz"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert" // исправлено: было stretchr -> stretchr
)

type AuthZServiceSuite struct {
	suite.Suite
}

func TestAuthZService(t *testing.T) {
	suite.RunSuite(t, new(AuthZServiceSuite))
}

func (s *AuthZServiceSuite) TestAuthZ_Authz_GetUserID(t provider.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ := authz.NewAuthZ()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.UserPayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.UserIDFromContext(ctx)
		sCtx.Require().NoError(err)
		assert.Equal(t, userID, resUserID)
	})
	t.WithNewStep("not found userID", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		_, err := authzServ.UserIDFromContext(ctx)

		assert.ErrorIs(t, err, authz.ErrNotAuthZ)
	})
}

func (s *AuthZServiceSuite) TestAuthZ_Authz_GetEmployeeID(t provider.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ := authz.NewAuthZ()
	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.EmployeePayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.EmployeeIDFromContext(ctx)
		sCtx.Require().NoError(err)
		assert.Equal(t, userID, resUserID)
	})
	t.WithNewStep("not found userID", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		_, err := authzServ.EmployeeIDFromContext(ctx)

		assert.ErrorIs(t, err, authz.ErrNotAuthZ)
	})
}

func (s *AuthZServiceSuite) TestAuthZ_Authz_GetAdminID(t provider.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ := authz.NewAuthZ()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.AdminPayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.AdminIDFromContext(ctx)
		sCtx.Require().NoError(err)
		assert.Equal(t, userID, resUserID)
	})
	t.WithNewStep("not found userID", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		_, err := authzServ.AdminIDFromContext(ctx)

		assert.ErrorIs(t, err, authz.ErrNotAuthZ)
	})
}
