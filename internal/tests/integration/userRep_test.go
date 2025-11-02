package integration_test

import (
	"context"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type UserRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx         context.Context
	userCreator testobj.UserMother
	userRep     userrep.UserRep
}

func TestUserRepSuite(t *testing.T) {
	suite.RunSuite(t, new(UserRepSuite))
}

func (s *UserRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.userCreator = testobj.NewUserMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.userRep, err = userrep.NewUserRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *UserRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "user")
}

func (s *UserRepSuite) AfterAll(t provider.T) {
	if s.userRep != nil {
		s.userRep.Close()
	}
}
func (s *UserRepSuite) TestUserRep_GetAll(t provider.T) {
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		users := []*models.User{
			s.userCreator.DefaultUserP(uuid.New()),
			s.userCreator.DefaultUserP(uuid.New()),
		}
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, users)
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, users)

		// ACT
		resUsers, err := s.userRep.GetAll(s.ctx)

		sCtx.Require().NoError(err)
		fixturesrep.AssertUsersAreInRes(t, users, resUsers)
	})
}

func (s *UserRepSuite) TestUserRep_GetAllSubscribed(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		subscribedUser := s.userCreator.UserWithSubscribeP(uuid.New(), true)
		unsubscribedUser := s.userCreator.UserWithSubscribeP(uuid.New(), false)

		users := []*models.User{subscribedUser, unsubscribedUser}
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, users)
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, users)

		// ACT
		resUsers, err := s.userRep.GetAllSubscribed(s.ctx)

		t.Require().NoError(err)
		fixturesrep.AssertUsersAreInRes(t, []*models.User{subscribedUser}, resUsers)
	})
}

func (s *UserRepSuite) TestUserRep_GetByID(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		users := []*models.User{
			s.userCreator.DefaultUserP(uuid.New()),
		}
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, users)
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, users)

		// ACT
		resUser, err := s.userRep.GetByID(s.ctx, users[0].GetID())

		t.Require().NoError(err)
		resUser.Equal(users[0])
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.userRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, userrep.ErrUserNotFound)
	})
}

func (s *UserRepSuite) TestUserRep_GetByLogin(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		users := []*models.User{
			s.userCreator.DefaultUserP(uuid.New()),
		}
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, users)
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, users)

		// ACT
		resUser, err := s.userRep.GetByLogin(s.ctx, users[0].GetLogin())

		t.Require().NoError(err)
		resUser.Equal(users[0])
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.userRep.GetByLogin(s.ctx, "nonexistent_login")

		t.Require().ErrorIs(err, userrep.ErrUserNotFound)
	})
}

func (s *UserRepSuite) TestUserRep_Add(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		user := s.userCreator.DefaultUserP(uuid.New())
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, []*models.User{user})

		// ACT
		err := s.userRep.Add(s.ctx, user)

		t.Require().NoError(err)

		// Verify user was added
		resUser, err := s.userRep.GetByID(s.ctx, user.GetID())
		t.Require().NoError(err)
		resUser.Equal(user)
	})

	t.Run("Duplicate login error", func(t provider.T) {
		user := s.userCreator.DefaultUserP(uuid.New())
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, []*models.User{user})
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, []*models.User{user})

		// ACT - try to add user with same login
		duplicateUser := s.userCreator.UserWithLoginP(uuid.New(), user.GetLogin())
		err := s.userRep.Add(s.ctx, duplicateUser)

		t.Require().ErrorIs(err, userrep.ErrDuplicateLoginUser)
	})
}

func (s *UserRepSuite) TestUserRep_Delete(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		user := s.userCreator.DefaultUserP(uuid.New())
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, []*models.User{user})

		// ACT
		err := s.userRep.Delete(s.ctx, user.GetID())

		t.Require().NoError(err)

		// Verify user was deleted
		_, err = s.userRep.GetByID(s.ctx, user.GetID())
		t.Require().ErrorIs(err, userrep.ErrUserNotFound)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.userRep.Delete(s.ctx, uuid.New())

		t.Require().ErrorIs(err, userrep.ErrRowsAffected)
	})
}

func (s *UserRepSuite) TestUserRep_Update(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		user := s.userCreator.DefaultUserP(uuid.New())
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, []*models.User{user})
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, []*models.User{user})

		newUser, err := models.NewUser(
			user.GetID(),
			"new_username",
			"new_login"+uuid.NewString(),
			user.GetHashedPassword(),
			user.GetCreatedAt(),
			"new_email@example.com",
			!user.IsSubscribedToMail(),
		)
		t.Require().NoError(err)

		funcUpdate := func(u *models.User) (*models.User, error) {
			return &newUser, nil
		}

		// ACT
		updatedUser, err := s.userRep.Update(s.ctx, user.GetID(), funcUpdate)

		t.Require().NoError(err)
		newUser.Equal(updatedUser)

		// Verify changes persisted
		dbUser, err := s.userRep.GetByID(s.ctx, user.GetID())
		t.Require().NoError(err)
		newUser.Equal(dbUser)
	})

	t.Run("Not found", func(t provider.T) {
		funcUpdate := func(u *models.User) (*models.User, error) {
			return u, nil
		}

		// ACT
		_, err := s.userRep.Update(s.ctx, uuid.New(), funcUpdate)

		t.Require().ErrorIs(err, userrep.ErrUserNotFound)
	})

	t.Run("Duplicate login error", func(t provider.T) {
		user1 := s.userCreator.DefaultUserP(uuid.New())
		user2 := s.userCreator.DefaultUserP(uuid.New())

		users := []*models.User{user1, user2}
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, users)
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, users)

		newUser, err := models.NewUser(
			user1.GetID(),
			user1.GetUsername(),
			user2.GetLogin(),
			user1.GetHashedPassword(),
			user1.GetCreatedAt(),
			user1.GetEmail(),
			user1.IsSubscribedToMail(),
		)
		t.Require().NoError(err)

		funcUpdate := func(u *models.User) (*models.User, error) {
			return &newUser, nil
		}

		// ACT
		_, err = s.userRep.Update(s.ctx, user1.GetID(), funcUpdate)

		t.Require().Error(err)
	})
}

func (s *UserRepSuite) TestUserRep_UpdateSubscribeToMailing(t provider.T) {
	t.Parallel()

	t.Run("Success - enable subscription", func(t provider.T) {
		user := s.userCreator.UserWithSubscribeP(uuid.New(), false)
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, []*models.User{user})
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, []*models.User{user})

		// ACT
		err := s.userRep.UpdateSubscribeToMailing(s.ctx, user.GetID(), true)

		t.Require().NoError(err)

		// Verify subscription updated
		updatedUser, err := s.userRep.GetByID(s.ctx, user.GetID())
		t.Require().NoError(err)
		t.Require().True(updatedUser.IsSubscribedToMail())
	})

	t.Run("Success - disable subscription", func(t provider.T) {
		user := s.userCreator.UserWithSubscribeP(uuid.New(), true)
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, []*models.User{user})
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, []*models.User{user})

		// ACT
		err := s.userRep.UpdateSubscribeToMailing(s.ctx, user.GetID(), false)

		t.Require().NoError(err)

		// Verify subscription updated
		updatedUser, err := s.userRep.GetByID(s.ctx, user.GetID())
		t.Require().NoError(err)
		t.Require().False(updatedUser.IsSubscribedToMail())
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.userRep.UpdateSubscribeToMailing(s.ctx, uuid.New(), true)

		t.Require().ErrorIs(err, userrep.ErrUserNotFound)
	})
}
