package integration_test

import (
	"context"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type AdminRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx          context.Context
	adminCreator testobj.AdminMother
	adminRep     adminrep.AdminRep
}

func TestAdminRepSuite(t *testing.T) {
	suite.RunSuite(t, new(AdminRepSuite))
}

func (s *AdminRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.adminCreator = testobj.NewAdminMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.adminRep, err = adminrep.NewAdminRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *AdminRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "admin")
}

func (s *AdminRepSuite) AfterAll(t provider.T) {
	if s.adminRep != nil {
		s.adminRep.Close()
	}
}

func (s *AdminRepSuite) TestAdminRep_GetAll(t provider.T) {
	t.Parallel()
	t.Run("Success", func(t provider.T) {
		admins := []*models.Admin{
			s.adminCreator.DefaultAdminP(uuid.New()),
			s.adminCreator.DefaultAdminP(uuid.New()),
		}
		fixturesrep.AddTestAdmin(t, s.ctx, s.adminRep, admins)
		defer fixturesrep.DelTestAdmin(t, s.ctx, s.adminRep, admins)

		// ACT
		resAdmins, err := s.adminRep.GetAll(s.ctx)

		t.Require().NoError(err)
		fixturesrep.AssertAdminsAreInRes(t, admins, resAdmins)
	})
}

func (s *AdminRepSuite) TestAdminRep_GetByID(t provider.T) {
	t.Parallel()
	t.Run("Success", func(t provider.T) {
		admins := []*models.Admin{
			s.adminCreator.DefaultAdminP(uuid.New()),
		}
		fixturesrep.AddTestAdmin(t, s.ctx, s.adminRep, admins)
		defer fixturesrep.DelTestAdmin(t, s.ctx, s.adminRep, admins)

		// ACT
		resAdmin, err := s.adminRep.GetByID(s.ctx, admins[0].GetID())

		t.Require().NoError(err)
		resAdmin.Equal(admins[0])
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.adminRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, adminrep.ErrAdminNotFound)
	})
}

func (s *AdminRepSuite) TestAdminRep_GetByLogin(t provider.T) {
	t.Parallel()
	t.Run("Success", func(t provider.T) {
		admins := []*models.Admin{
			s.adminCreator.DefaultAdminP(uuid.New()),
		}
		fixturesrep.AddTestAdmin(t, s.ctx, s.adminRep, admins)
		defer fixturesrep.DelTestAdmin(t, s.ctx, s.adminRep, admins)

		// ACT
		resAdmin, err := s.adminRep.GetByLogin(s.ctx, admins[0].GetLogin())

		t.Require().NoError(err)
		resAdmin.Equal(admins[0])
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.adminRep.GetByLogin(s.ctx, "---")

		t.Require().ErrorIs(err, adminrep.ErrAdminNotFound)
	})
}

func (s *AdminRepSuite) TestAdminRep_Add(t provider.T) {
	t.Parallel()
	t.Run("Duplicate login error", func(t provider.T) {
		admin := s.adminCreator.DefaultAdminP(uuid.New())
		fixturesrep.AddTestAdmin(t, s.ctx, s.adminRep, []*models.Admin{admin})
		defer fixturesrep.DelTestAdmin(t, s.ctx, s.adminRep, []*models.Admin{admin})

		// ACT
		err := s.adminRep.Add(s.ctx, admin)

		t.Require().ErrorIs(err, adminrep.ErrDuplicateLoginAdm)
	})
}

func (s *AdminRepSuite) TestAdminRep_Update(t provider.T) {
	t.Parallel()
	t.Run("Duplicate login error", func(t provider.T) {
		admin := s.adminCreator.DefaultAdminP(uuid.New())
		fixturesrep.AddTestAdmin(t, s.ctx, s.adminRep, []*models.Admin{admin})
		defer fixturesrep.DelTestAdmin(t, s.ctx, s.adminRep, []*models.Admin{admin})

		newAdmin, err := models.NewAdmin(
			admin.GetID(),
			"new_username",
			"new_login"+uuid.NewString(),
			"new_hashed_password",
			admin.GetCreatedAt(),
			false,
		)
		t.Require().NoError(err)
		funcUpdate := func(a *models.Admin) (*models.Admin, error) {
			return &newAdmin, err
		}

		// ACT
		err = s.adminRep.Update(s.ctx, admin.GetID(), funcUpdate)

		t.Require().NoError(err)
		unpdatedAdmin, err := s.adminRep.GetByID(s.ctx, admin.GetID())
		t.Require().NoError(err)
		newAdmin.Equal(unpdatedAdmin)
	})
}
