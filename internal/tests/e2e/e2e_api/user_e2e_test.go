package e2eapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/otp"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailreader"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type UserE2ESuite struct {
	fixtures.BaseE2ESuite
	client                 *fixtures.HTTPClient
	registerUserReqCreator testobj.AuthUserRequestMother
	registerUserData       authmodels.RegisterUserRequest
}

func TestUserE2E(t *testing.T) {
	suite.RunSuite(t, new(UserE2ESuite))

}

func (s *UserE2ESuite) BeforeAll(t provider.T) {
	s.BaseE2ESuite.BeforeAll(t)
	s.client = fixtures.NewHTTPClient(s.BaseURL)
}

func (s *UserE2ESuite) BeforeEach(t provider.T) {
	s.registerUserReqCreator = testobj.NewRegisterUserRequestMother()
	s.registerUserData = s.registerUserReqCreator.RegisterWithEmail(s.EmailReader.Username)
	t.WithNewStep("Register new user", func(sCtx provider.StepCtx) {
		resp, err := s.client.DoRequest("POST", "/auth-user/register", s.registerUserData)
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(http.StatusOK, resp.StatusCode)
		// respText := ""
		// bodyBytes, err := io.ReadAll(resp.Body)
		// if err == nil {
		// 	respText = string(bodyBytes)
		// }
		// fmt.Printf("respText: %s\n\n", respText)

	})
}

func (s *UserE2ESuite) AfterEach(t provider.T) {
	t.WithNewStep("Delete user", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		userRep, err := userrep.NewUserRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		user, err := userRep.GetByLogin(ctx, s.registerUserData.Login)
		sCtx.Require().NoError(err)
		userRep.Delete(ctx, user.GetID())
	})
}

func (s *UserE2ESuite) TestUser_MVP(t provider.T) {
	t.Tag("e2e")
	t.Description("User register, auth, search events")

	// var accessToken string

	t.WithNewStep("Login new user", func(sCtx provider.StepCtx) {
		loginReq := s.registerUserReqCreator.Login(s.registerUserData.Login, s.registerUserData.Password)
		resp, err := s.client.DoRequest("POST", "/auth-user/login", loginReq)

		sCtx.Require().NoError(err)
		sCtx.Require().Equal(http.StatusOK, resp.StatusCode)

		var responseBody authmodels.LoginUserResponse
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		sCtx.Require().NoError(err)

		sessionID := responseBody.SessionAuthID
		sCtx.Require().NotEmpty(sessionID)
		_, err = uuid.Parse(responseBody.SessionAuthID)
		sCtx.Require().NoError(err)

		time.Sleep(time.Second)
		email, err := s.EmailReader.FindEmailByCriteria(emailreader.SearchCriteria{
			From:    s.EmailCnfg.From,
			Subject: otp.OTPSubject,
		})
		sCtx.Require().NoError(err)
		OTPCode := strings.ReplaceAll(email.Body, "\r\n", "")

		login2faReq := authmodels.Login2FAUserRequest{
			SessionAuthID: sessionID,
			OTPCode:       OTPCode,
		}

		resp, err = s.client.DoRequest("POST", "/auth-user/login-2fa", login2faReq)
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(http.StatusOK, resp.StatusCode)
		// respText := ""
		// bodyBytes, err := io.ReadAll(resp.Body)
		// if err == nil {
		// 	respText = string(bodyBytes)
		// }
		// fmt.Printf("respText: %s\n\n", respText)

		var loginResp authmodels.Login2FAUserResponse
		err = s.client.GetBody(resp, &loginResp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(loginResp.AccessToken)
		// accessToken := loginResp.AccessToken
	})

	eventCreator := testobj.NewEventMother()
	events := []*models.Event{
		eventCreator.EventP(uuid.New()),
		eventCreator.EventP(uuid.New()),
		eventCreator.EventP(uuid.New()),
	}
	t.WithNewStep("Add events", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventRep, err := eventrep.NewEventRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		employeeRep, err := employeerep.NewEmployeeRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		adminRep, err := adminrep.NewAdminRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		artworkRep, err := artworkrep.NewArtworkRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		authorRep, err := authorrep.NewAuthorRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		collectionRep, err := collectionrep.NewCollectionRep(ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		fixturesrep.AddTestEvents(t, ctx, events, eventRep, employeeRep, adminRep, artworkRep, authorRep, collectionRep)
	})

	t.WithNewStep("Search events", func(sCtx provider.StepCtx) {
		resp, err := s.client.DoRequest("GET", "/museum/events", nil)

		sCtx.Require().NoError(err)
		sCtx.Require().Equal(http.StatusOK, resp.StatusCode)
		var eventResp []jsonreqresp.EventResponse
		err = s.client.GetBody(resp, &eventResp)
		sCtx.Require().NoError(err)
		expectedEventResp := make([]jsonreqresp.EventResponse, len(events))
		for i, v := range events {
			expectedEventResp[i] = v.ToEventResponse()
		}
		fixturesrep.AssertEventResponsesAreInRes(t, eventResp, expectedEventResp)
	})

}
