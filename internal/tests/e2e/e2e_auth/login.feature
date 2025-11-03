Feature: User login
  As an authenticated user
  I want to login 
  So that I can access my personal account

  Scenario: Successful user login
    Given I have valid login data
    When I send POST request to "/auth-user/login" with valid data
    Then the response status should be 200
    And the response must contain the session ID
    Given I have valid otp code
    When I send POST request to "/auth-user/login-2fa" with valid code
    Then the response status should be 200
    And the response must contain the access token

  Scenario: Limited number of input attempts
    Given I have invalid login data
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 401
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 401
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 401
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 429

  Scenario: Recovery when the attempt limit is exceeded
    Given I have invalid login data
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 401
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 401
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 401
    When I send POST request to "/auth-user/login" with invalid data
    Then the response status should be 429
    And wait duration session
    Given I have valid login data
    When I send POST request to "/auth-user/login" with valid data
    Then the response status should be 200
    And the response must contain the session ID

  Scenario: Change password
    Given I have valid login data
    When I send POST request to "/auth-user/login" with valid data
    Then the response status should be 200
    And the response must contain the session ID
    Given I have valid otp code
    When I send POST request to "/auth-user/login-2fa" with valid code
    Then the response status should be 200
    And the response must contain the access token

    And I have changed password
    Then the response status should be 200

    Given I have valid login data
    When I send POST request to "/auth-user/login" with valid data
    Then the response status should be 200
    And the response must contain the session ID
    Given I have valid otp code
    When I send POST request to "/auth-user/login-2fa" with valid code
    Then the response status should be 200
    And the response must contain the access token