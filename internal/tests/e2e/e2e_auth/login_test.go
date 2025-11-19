package e2eauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/otp"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailreader"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

func TestLoginFeatures(t *testing.T) {
	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeLoginScenario,
		Options: &godog.Options{
			TestingT: t,
			Format:   "pretty",
			Paths:    []string{"./login.feature"}, // Поиск .feature файлов
			Strict:   true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run login feature tests")
	}
}

func InitializeLoginScenario(sc *godog.ScenarioContext) {
	regCtx := &userLoginContext{}
	var userRep userrep.UserRep = nil
	var err error

	sc.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
		regCtx.baseConf, err = NewGodogTestConfig()
		if err != nil {
			return ctx, err
		}
		regCtx.client = fixtures.NewHTTPClient(regCtx.baseConf.BaseURL)

		userRep, err = userrep.NewUserRep(ctx,
			regCtx.baseConf.AppCnfg.Datebase, regCtx.baseConf.DBCreds, regCtx.baseConf.DBCnfg)
		if err != nil {
			return ctx, err
		}

		registerUserReqCreator := testobj.NewRegisterUserRequestMother()
		// "tmpforread@mail.ru"
		registerReq := registerUserReqCreator.RegisterWithEmail(regCtx.baseConf.EmailReader.Username)
		resp, err := regCtx.client.DoRequest("POST", "/auth-user/register", registerReq)
		if err != nil {
			return ctx, err
		}
		if resp.StatusCode != 200 {
			respBody := ""
			if resp.Body != nil {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err == nil {
					respBody = string(bodyBytes)
				}
			}
			return ctx, fmt.Errorf("got status %d (%s)", resp.StatusCode, respBody)
		}
		regCtx.registerData = &registerReq

		return ctx, nil
	})

	sc.Step(`^I have valid login data$`, regCtx.iHaveValidLoginData)
	sc.Step(`^I have invalid login data$`, regCtx.iHaveInvalidLoginData)
	sc.Step(`^I send POST request to "([^"]*)" with (valid data|duplicate data|invalid data)$`,
		regCtx.iSendPOSTRequestToWithData)
	sc.Step(`^the response status should be (\d+)$`, regCtx.theResponseStatusShouldBe)
	sc.Step(`^the response must contain the session ID$`, regCtx.theResponseMustContainTheSessionID)
	sc.Step(`^I have valid otp code$`, regCtx.iHaveValidOTPCode)
	sc.Step(`^I send POST request to "([^"]*)" with (valid code|invalid code)$`,
		regCtx.iSendPOSTRequestToWithCode)
	sc.Step(`^the response must contain the access token$`,
		regCtx.theResponseMustContainTheAccessToken)
	sc.Step(`^wait duration session$`, regCtx.wait)
	sc.Step(`^I have changed password$`, regCtx.iHaveChangedPassword)

	sc.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			fmt.Printf("❌ Scenario failed: %s - %v\n", scenario.Name, err)
		} else {
			fmt.Printf("✅ Scenario passed: %s\n", scenario.Name)
		}

		if userRep != nil {
			user, err := userRep.GetByLogin(ctx, regCtx.registerData.Login)
			if err == nil {
				_ = userRep.Delete(ctx, user.GetID())
			}
		}

		return ctx, nil
	})
}

type userLoginContext struct {
	baseConf     *GodogTestConfig
	client       *fixtures.HTTPClient
	registerData *authmodels.RegisterUserRequest
	loginReq     *authmodels.LoginUserRequest
	response     *http.Response
	sessionID    string
	otpCode      string
	accessToken  string
}

func (r *userLoginContext) iHaveValidLoginData() error {
	r.loginReq = &authmodels.LoginUserRequest{
		Login:    r.registerData.Login,
		Password: r.registerData.Password,
	}
	return nil
}

func (r *userLoginContext) iHaveInvalidLoginData() error {
	r.loginReq = &authmodels.LoginUserRequest{
		Login:    r.registerData.Login,
		Password: "aaaa",
	}
	return nil
}

func (r *userLoginContext) iSendPOSTRequestToWithData(path string, dataType string) error {
	if !(dataType == "valid data" || dataType == "duplicate data" || dataType == "invalid data") {
		return fmt.Errorf("unknown data type: %s", dataType)
	}
	resp, err := r.client.DoRequest("POST", path, *r.loginReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	r.response = resp
	return nil
}

func (r *userLoginContext) theResponseStatusShouldBe(status int) error {
	if r.response.StatusCode != status {
		resp := ""
		if r.response.Body != nil {
			bodyBytes, err := io.ReadAll(r.response.Body)
			if err == nil {
				resp = string(bodyBytes)
			}
		}
		return fmt.Errorf("expected status %d, but got %d (%s)", status, r.response.StatusCode, resp)
	}
	return nil
}

func (r *userLoginContext) theResponseMustContainTheSessionID() error {
	if r.response.Body == nil {
		return fmt.Errorf("expected not empty response body")
	}
	defer r.response.Body.Close()
	var responseBody authmodels.LoginUserResponse
	if err := json.NewDecoder(r.response.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}
	if responseBody.SessionAuthID == "" {
		return fmt.Errorf("expected sessionAuthID in response, but got empty")
	}
	if _, err := uuid.Parse(responseBody.SessionAuthID); err != nil {
		return fmt.Errorf("sessionAuthID is not a valid UUID: %s, error: %v",
			responseBody.SessionAuthID, err)
	}
	r.sessionID = responseBody.SessionAuthID
	return nil
}

func (r *userLoginContext) iHaveValidOTPCode() error {
	time.Sleep(time.Second)
	emails, err := r.baseConf.EmailReader.FindEmailByCriteria(emailreader.SearchCriteria{
		From:    r.baseConf.EmailCnfg.From,
		Subject: otp.OTPSubject,
	})
	if err != nil {
		return fmt.Errorf("Error read emails: %w", err)
	}
	r.otpCode = strings.ReplaceAll(emails.Body, "\r\n", "")
	return nil
}

func (r *userLoginContext) iSendPOSTRequestToWithCode(path string, codeType string) error {
	if !(codeType == "valid code" || codeType == "invalid code") {
		return fmt.Errorf("unknown code type: %s", codeType)
	}
	req := authmodels.Login2FAUserRequest{
		SessionAuthID: r.sessionID,
		OTPCode:       r.otpCode,
	}
	resp, err := r.client.DoRequest("POST", path, req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	r.response = resp
	return nil
}

func (r *userLoginContext) theResponseMustContainTheAccessToken() error {
	if r.response.Body == nil {
		return fmt.Errorf("expected not empty response body")
	}
	defer r.response.Body.Close()
	var responseBody authmodels.Login2FAUserResponse
	if err := json.NewDecoder(r.response.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}
	if responseBody.AccessToken == "" {
		return fmt.Errorf("expected AccessToken in response, but got empty")
	}
	r.accessToken = responseBody.AccessToken
	return nil
}

func (r *userLoginContext) wait() error {
	seconds := r.baseConf.AppCnfg.DurationLoginSession.Seconds() + 1
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}

func (r *userLoginContext) iHaveChangedPassword() error {
	r.registerData.Password = uuid.NewString()[:10]
	req := jsonreqresp.ChangePasswordRequest{
		Password: r.registerData.Password,
	}
	resp, err := r.client.DoAuthRequest("PUT", "/user/self/change-password", req, r.accessToken)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	r.response = resp
	return nil
}
