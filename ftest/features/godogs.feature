# file: $GOPATH/godogs/features/godogs.feature
Feature: go-test
  In order to get the coverage report

  Scenario: With Empty Data
    When I send a "GET" request to "http://localhost:8085/v1/todo"
    Then the response code should be 404
