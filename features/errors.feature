Feature: Error Handling
  As an API consumer
  I want consistent error responses
  So that I can handle errors programmatically

  Scenario: Accessing undefined route returns 404
    Given the server is running
    When I send a GET request to "/undefined-route"
    Then the response status code should be 404
    And the response content type should be "application/json"
    And the response body should contain "error"
    And the response body should contain "not found"

  Scenario: Accessing undefined route with POST returns 404
    Given the server is running
    When I send a POST request to "/undefined-route"
    Then the response status code should be 404
    And the response content type should be "application/json"
