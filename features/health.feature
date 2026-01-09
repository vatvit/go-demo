Feature: Health Check Endpoint
  As a system operator
  I want a health check endpoint
  So that I can monitor if the server is running

  Scenario: Server returns healthy status
    Given the server is running
    When I send a GET request to "/health"
    Then the response status code should be 200
    And the response content type should be "application/json"
    And the response body should contain "status"
    And the response body should contain "ok"
