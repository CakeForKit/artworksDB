package e2eauth

// type registrationContext struct {
// 	client          *fixtures.HTTPClient
// 	baseURL         string
// 	registerRequest *authmodels.RegisterUserRequest
// 	response        *http.Response
// 	responseBody    map[string]interface{}
// }

// func (r *registrationContext) iHaveValidRegistrationData() error {
// 	// Генерируем уникальные данные для регистрации
// 	uniqueID := uuid.New().String()[:8]
// 	r.registerRequest = &authmodels.RegisterUserRequest{
// 		Login:          fmt.Sprintf("testuser_%s", uniqueID),
// 		Password:       "securePassword123",
// 		Email:          fmt.Sprintf("test_%s@example.com", uniqueID),
// 		SubscribeEmail: true,
// 		Username:       "User",
// 	}
// 	return nil
// }

// func (r *registrationContext) iHaveRegistrationDataWithDuplicateLogin() error {
// 	// Используем существующий логин (предполагая, что пользователь уже зарегистрирован)
// 	r.registerRequest = &authmodels.RegisterUserRequest{
// 		Login:          "existing_user", // Должен существовать в базе
// 		Password:       "securePassword123",
// 		Email:          "existing@example.com",
// 		SubscribeEmail: true,
// 		Username:       "User",
// 	}
// 	return nil
// }

// func (r *registrationContext) iHaveInvalidRegistrationData() error {
// 	// Создаем невалидные данные
// 	r.registerRequest = &authmodels.RegisterUserRequest{
// 		Login:          "",              // Пустой логин
// 		Password:       "short",         // Слишком короткий пароль
// 		Email:          "invalid-email", // Невалидный email
// 		SubscribeEmail: false,           // Пустое имя
// 		Username:       "",              // Пустая фамилия
// 	}
// 	return nil
// }

// func (r *registrationContext) iSendPOSTRequestToWithData(path string, dataType string) error {
// 	var reqBody []byte
// 	var err error

// 	switch dataType {
// 	case "valid data":
// 		reqBody, err = json.Marshal(r.registerRequest)
// 	case "duplicate data":
// 		reqBody, err = json.Marshal(r.registerRequest)
// 	case "invalid data":
// 		reqBody, err = json.Marshal(r.registerRequest)
// 	default:
// 		return fmt.Errorf("unknown data type: %s", dataType)
// 	}

// 	if err != nil {
// 		return fmt.Errorf("failed to marshal request: %w", err)
// 	}

// 	resp, err := r.client.DoRequest("POST", path, bytes.NewReader(reqBody))
// 	if err != nil {
// 		return fmt.Errorf("failed to send request: %w", err)
// 	}

// 	r.response = resp

// 	// Парсим тело ответа, если есть
// 	if resp.Body != nil && resp.StatusCode != http.StatusOK {
// 		defer resp.Body.Close()
// 		r.responseBody = make(map[string]interface{})
// 		if err := json.NewDecoder(resp.Body).Decode(&r.responseBody); err != nil {
// 			// Игнорируем ошибки парсинга, так как ответ может быть пустым
// 			fmt.Printf("Warning: failed to parse response body: %v\n", err)
// 		}
// 	}

// 	return nil
// }

// func (r *registrationContext) theResponseStatusShouldBe(status int) error {
// 	if r.response.StatusCode != status {
// 		return fmt.Errorf("expected status %d, but got %d", status, r.response.StatusCode)
// 	}
// 	return nil
// }

// func (r *registrationContext) theResponseShouldBeEmpty() error {
// 	// Проверяем, что тело ответа пустое или содержит только {}
// 	if r.response.Body != nil {
// 		defer r.response.Body.Close()
// 		var body map[string]interface{}
// 		if err := json.NewDecoder(r.response.Body).Decode(&body); err == nil {
// 			if len(body) > 0 {
// 				return fmt.Errorf("expected empty response, but got: %v", body)
// 			}
// 		}
// 	}
// 	return nil
// }

// func (r *registrationContext) theResponseShouldContainErrorAboutDuplicateLogin() error {
// 	if r.responseBody == nil {
// 		return fmt.Errorf("expected error response, but got empty body")
// 	}

// 	errorMsg, exists := r.responseBody["error"].(string)
// 	if !exists {
// 		return fmt.Errorf("expected 'error' field in response")
// 	}

// 	// Проверяем, что ошибка связана с дубликатом
// 	if errorMsg == "" {
// 		return fmt.Errorf("expected error message about duplicate login")
// 	}

// 	fmt.Printf("Got error message: %s\n", errorMsg)
// 	return nil
// }

// func (r *registrationContext) theResponseShouldContainValidationError() error {
// 	if r.responseBody == nil {
// 		return fmt.Errorf("expected validation error response, but got empty body")
// 	}

// 	_, exists := r.responseBody["error"].(string)
// 	if !exists {
// 		return fmt.Errorf("expected 'error' field in validation response")
// 	}

// 	return nil
// }
