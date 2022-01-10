Feature: Behaviour of application when performing requests against /v1/policies endpoints - testing the authorisation middleware functionality

      # A Background applies to all scenarios in this Feature
  Background:
    Given I have these policies:
      """
      [
          {
              "id": "admin",
              "role": "admin",
              "entities": [
                "group/admin"
              ],
              "conditions": []
          },
          {
              "id": "publisher",
              "role": "publisher",
              "entities": [
                "group/publisher"
              ],
              "conditions": []
          },
          {
              "id": "viewer",
              "role": "viewer",
              "entities": [
                "group/viewer"
              ],
              "conditions": [
                  {
                      "operator": "StringEquals",
                      "attributes": [
                        "collection-id"
                      ],
                      "values": [
                        "collection-765"
                      ]
                  }
              ]
          }
      ]
      """

  Scenario: [Test #1] GET /v1/policies/publisher with a valid JWT token the response status is 200
    Given I am a publisher user
    When I GET "/v1/policies/publisher"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
      """
      {
              "id": "publisher",
              "role": "publisher",
              "entities": [
                "group/publisher"
              ]
          }
      """

  Scenario: [Test #2] GET /v1/policies/publisher with invalid JWT token in header - the response status is 403 (forbidden)
    Given I am a publisher user with invalid auth token
    When I GET "/v1/policies/publisher"
    Then the HTTP status code should be "403"

  Scenario: [Test #3] GET /v1/policies/viewer to fetch a policy having all parameters
    Given I am a viewer user
    When I GET "/v1/policies/viewer"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
      """
      {
              "id": "viewer",
              "role": "viewer",
              "entities": [
                "group/viewer"
              ],
              "conditions": [
                  {
                      "operator": "StringEquals",
                      "attributes": [
                        "collection-id"
                      ],
                      "values": [
                        "collection-765"
                      ]
                  }
              ]
          }
      """

  Scenario: [Test #4] GET /v1/policies/admin with invalid JWT token in header - the response status is 403 (forbidden)
    Given I am a basic user
    When I GET "/v1/policies/admin"
    Then the HTTP status code should be "403"

  Scenario: [Test #5] Receive not found when doing a GET for a non existent policy
    Given I am a publisher user
    And I GET "/v1/policies/notFound"
    Then the HTTP status code should be "404"
    And the response header "Content-Type" should be "text/plain; charset=utf-8"
    And I should receive the following response:
      """
      policy not found
      """"""
