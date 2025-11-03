package e2eauth

// func TestLoginFeatures1(t *testing.T) {
// 	suite := godog.TestSuite{
// 		ScenarioInitializer: InitializeLoginScenario,
// 		Options: &godog.Options{
// 			TestingT: t,
// 			Format:   "pretty",
// 			Paths:    []string{"internal/tests/e2e/e2e_auth/"}, // Поиск .feature файлов
// 			Strict:   true,
// 		},
// 	}

// 	if suite.Run() != 0 {
// 		t.Fatal("non-zero status returned, failed to run registration feature tests")
// 	}
// }

// func InitializeTestSuite1(ctx *godog.TestSuiteContext) {
// 	// ВЫПОЛНЯЕТСЯ ОДИН РАЗ для всей Test Suite
// 	ctx.BeforeSuite(func() {
// 		// ◀️ Аналог BeforeAll в ozonteach
// 		fmt.Println("🚀 Starting Test Suite - выполняется ОДИН РАЗ перед всеми фичами")

// 		// Здесь обычно:
// 		// - Запуск приложения/сервера
// 		// - Миграции базы данных
// 		// - Инициализация глобальных ресурсов
// 		// - Подготовка тестового окружения
// 	})

// 	ctx.AfterSuite(func() {
// 		// ◀️ Аналог AfterAll в ozonteach
// 		fmt.Println("✅ Test Suite Completed - выполняется ОДИН РАЗ после всех фич")

// 		// Здесь обычно:
// 		// - Остановка приложения/сервера
// 		// - Очистка базы данных
// 		// - Освобождение ресурсов
// 		// - Генерация отчетов
// 	})
// }

// func InitializeScenario(sc *godog.ScenarioContext) {
// 	// ВЫПОЛНЯЕТСЯ для КАЖДОГО сценария
// 	sc.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
// 		// ◀️ Аналог BeforeEach в ozonteach
// 		fmt.Printf("📝 Starting scenario: %s - выполняется ПЕРЕД КАЖДЫМ сценарием\n", scenario.Name)

// 		// Здесь обычно:
// 		// - Сброс состояния тестового контекста
// 		// - Очистка данных пользователя
// 		// - Инициализация клиента HTTP
// 		// - Подготовка свежих данных для теста

// 		return ctx, nil
// 	})

// 	sc.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
// 		// ◀️ Аналог AfterEach в ozonteach
// 		if err != nil {
// 			fmt.Printf("❌ Scenario failed: %s - %v\n", scenario.Name, err)
// 		} else {
// 			fmt.Printf("✅ Scenario passed: %s\n", scenario.Name)
// 		}

// 		// Здесь обычно:
// 		// - Логирование результатов
// 		// - Снятие скриншотов (для UI тестов)
// 		// - Очистка временных данных
// 		// - Закрытие соединений

// 		return ctx, nil
// 	})
// }

// func InitializeLoginScenario1(sc *godog.ScenarioContext) {
// 	regCtx := &registrationContext{}

// 	// Before каждого сценария
// 	sc.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
// 		fmt.Printf("📝 Starting scenario: %s\n", scenario.Name)

// 		// Инициализация контекста
// 		regCtx.baseURL = "http://localhost:8080"
// 		regCtx.client = fixtures.NewHTTPClient(regCtx.baseURL)
// 		regCtx.registerRequest = nil
// 		regCtx.response = nil
// 		regCtx.responseBody = nil

// 		return ctx, nil
// 	})

// 	// sc.Step(`^I have valid registration data$`, regCtx.iHaveValidRegistrationData)
// 	// sc.Step(`^I have registration data with duplicate login$`, regCtx.iHaveRegistrationDataWithDuplicateLogin)
// 	// sc.Step(`^I have invalid registration data$`, regCtx.iHaveInvalidRegistrationData)
// 	// sc.Step(`^I send POST request to "([^"]*)" with (valid data|duplicate data|invalid data)$`, regCtx.iSendPOSTRequestToWithData)
// 	// sc.Step(`^the response status should be (\d+)$`, regCtx.theResponseStatusShouldBe)
// 	// sc.Step(`^the response should be empty$`, regCtx.theResponseShouldBeEmpty)
// 	// sc.Step(`^the response should contain error about duplicate login$`, regCtx.theResponseShouldContainErrorAboutDuplicateLogin)
// 	// sc.Step(`^the response should contain validation error$`, regCtx.theResponseShouldContainValidationError)

// 	// After каждого сценария
// 	sc.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
// 		if err != nil {
// 			fmt.Printf("❌ Scenario failed: %s - %v\n", scenario.Name, err)
// 		} else {
// 			fmt.Printf("✅ Scenario passed: %s\n", scenario.Name)
// 		}
// 		return ctx, nil
// 	})
// }

// type userLoginContext1 struct {
// 	baseSuite     *fixtures.BaseE2ESuite
// 	client        *fixtures.HTTPClient
// 	registerReq   *authmodels.RegisterUserRequest
// 	loginReq      *authmodels.LoginUserRequest
// 	loginResp     *authmodels.LoginUserResponse
// 	events        []*models.Event
// 	eventResponse []interface{}
// 	accessToken   string
// 	lastResponse  *http.Response
// }

// func (u *userLoginContext) theApplicationIsConfiguredAndRunning(ctx context.Context) (context.Context, error) {
// 	// Инициализация базовой конфигурации (аналог BeforeAll)
// 	if u.baseSuite == nil {
// 		u.baseSuite = &fixtures.BaseE2ESuite{}

// 		// Создаем mock provider.T для инициализации
// 		mockT := &mockProviderT{}
// 		u.baseSuite.BeforeAll(mockT)

// 		if u.baseSuite.BaseURL == "" {
// 			return ctx, fmt.Errorf("failed to initialize application configuration")
// 		}

// 		u.client = fixtures.NewHTTPClient(u.baseSuite.BaseURL)
// 	}
// 	return ctx, nil
// }
