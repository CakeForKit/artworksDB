package auth_test

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stateio/testify/require"
	"github.com/stretchr/testify/assert" // исправлено: было stateio -> stretchr
)

func TestAuthZ_Authz_GetUserID(t *testing.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ, err := auth.NewAuthZ()
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.UserPayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.UserIDFromContext(ctx)
		require.Nil(t, err)
		assert.Equal(t, userID, resUserID)
	})
	t.Run("not found userID", func(t *testing.T) {
		ctx := context.Background()

		_, err := authzServ.UserIDFromContext(ctx)

		assert.ErrorIs(t, err, auth.ErrNotAuthZ)
	})
}

func TestAuthZ_Authz_GetEmployeeID(t *testing.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ, err := auth.NewAuthZ()
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.EmployeePayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.EmployeeIDFromContext(ctx)
		require.Nil(t, err)
		assert.Equal(t, userID, resUserID)
	})
	t.Run("not found userID", func(t *testing.T) {
		ctx := context.Background()

		_, err := authzServ.EmployeeIDFromContext(ctx)

		assert.ErrorIs(t, err, auth.ErrNotAuthZ)
	})
}

func TestAuthZ_Authz_GetAdminID(t *testing.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ, err := auth.NewAuthZ()
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.AdminPayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.AdminIDFromContext(ctx)
		require.Nil(t, err)
		assert.Equal(t, userID, resUserID)
	})
	t.Run("not found userID", func(t *testing.T) {
		ctx := context.Background()

		_, err := authzServ.AdminIDFromContext(ctx)

		assert.ErrorIs(t, err, auth.ErrNotAuthZ)
	})
}
