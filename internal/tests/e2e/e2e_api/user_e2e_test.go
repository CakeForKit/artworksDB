package e2eapi_test

// import (
// 	"context"
// 	"net/http"
// 	"testing"

// 	"github.com/CakeForKit/artworksDB.git/internal/models"
// 	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
// 	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
// 	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
// 	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
// 	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
// 	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
// 	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
// 	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
// 	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
// 	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
// 	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
// 	"github.com/google/uuid"
// 	"github.com/ozontech/allure-go/pkg/framework/provider"
// 	"github.com/ozontech/allure-go/pkg/framework/suite"
// )

// type UserE2ESuite struct {
// 	fixtures.BaseE2ESuite
// 	client *fixtures.HTTPClient
// }

// func TestUserE2E(t *testing.T) {
// 	suite.RunSuite(t, new(UserE2ESuite))

// }

// func (s *UserE2ESuite) BeforeAll(t provider.T) {
// 	s.BaseE2ESuite.BeforeAll(t)
// 	s.client = fixtures.NewHTTPClient(s.BaseURL)
// }

// func (s *UserE2ESuite) TestUser_MVP(t provider.T) {
// 	t.Tag("e2e")
// 	t.Description("User register, auth, search events")

// 	// var accessToken string
// 	registerUserReqCreator := testobj.NewRegisterUserRequestMother()
// 	registerReq := registerUserReqCreator.RegisterDefault()
// 	t.WithNewStep("Register new user", func(sCtx provider.StepCtx) {
// 		resp, err := s.client.DoRequest("POST", "/auth-user/register", registerReq)
// 		sCtx.Require().NoError(err)
// 		sCtx.Require().Equal(http.StatusOK, resp.StatusCode)
// 	})

// 	t.WithNewStep("Login new user", func(sCtx provider.StepCtx) {
// 		loginReq := registerUserReqCreator.Login(registerReq.Login, registerReq.Password)
// 		resp, err := s.client.DoRequest("POST", "/auth-user/login", loginReq)

// 		sCtx.Require().NoError(err)
// 		sCtx.Require().Equal(http.StatusOK, resp.StatusCode)
// 		var loginResp authmodels.LoginUserResponse
// 		err = s.client.GetBody(resp, &loginResp)
// 		sCtx.Require().NoError(err)
// 		sCtx.Require().NotEmpty(loginResp.AccessToken)
// 		// accessToken := loginResp.AccessToken
// 	})

// 	eventCreator := testobj.NewEventMother()
// 	events := []*models.Event{
// 		eventCreator.EventP(uuid.New()),
// 		eventCreator.EventP(uuid.New()),
// 		eventCreator.EventP(uuid.New()),
// 	}
// 	t.WithNewStep("Add events", func(sCtx provider.StepCtx) {
// 		ctx := context.Background()
// 		eventRep, err := eventrep.NewEventRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
// 		sCtx.Require().NoError(err)
// 		employeeRep, err := employeerep.NewEmployeeRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
// 		sCtx.Require().NoError(err)
// 		adminRep, err := adminrep.NewAdminRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
// 		sCtx.Require().NoError(err)
// 		artworkRep, err := artworkrep.NewArtworkRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
// 		sCtx.Require().NoError(err)
// 		authorRep, err := authorrep.NewAuthorRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
// 		sCtx.Require().NoError(err)
// 		collectionRep, err := collectionrep.NewCollectionRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
// 		sCtx.Require().NoError(err)
// 		fixturesrep.AddTestEvents(t, ctx, events, eventRep, employeeRep, adminRep, artworkRep, authorRep, collectionRep)
// 	})

// 	t.WithNewStep("Search events", func(sCtx provider.StepCtx) {
// 		resp, err := s.client.DoRequest("GET", "/museum/events", nil)

// 		sCtx.Require().NoError(err)
// 		sCtx.Require().Equal(http.StatusOK, resp.StatusCode)
// 		var eventResp []jsonreqresp.EventResponse
// 		err = s.client.GetBody(resp, &eventResp)
// 		sCtx.Require().NoError(err)
// 		expectedEventResp := make([]jsonreqresp.EventResponse, len(events))
// 		for i, v := range events {
// 			expectedEventResp[i] = v.ToEventResponse()
// 		}
// 		fixturesrep.AssertEventResponsesAreInRes(t, eventResp, expectedEventResp)
// 	})
// }
