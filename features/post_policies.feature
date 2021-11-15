Feature: Behaviour of application when doing the POST /v1/policies endpoint, using a stripped down version of the database

  Scenario: [Test #1] POST /v1/policies with all the parameters
    Given I am an admin user
    When I POST "/v1/policies"
      """
      {
          "entities": [
            "e1",
            "e2"
          ],
          "role": "r1",
          "conditions": [
            {
              "attributes": [
                "a1"
              ],
              "operator": "and",
              "values": [
                "v1"
              ]
            }
          ]
      }
      """
    Then the HTTP status code should be "201"


  Scenario: [Test #2] POST /v1/policies without conditions
    Given I am an admin user
    When I POST "/v1/policies"
      """
      {
          "entities": [
            "e1",
            "e2"
          ],
          "role": "r1"
      }
      """
    Then the HTTP status code should be "201"

  Scenario: [Test #3] POST /v1/policies without entities
    Given I am an admin user
    When I POST "/v1/policies"
      """
      {
          "role": "r1",
          "conditions": [
            {
              "attributes": [
                "a1"
              ],
              "operator": "and",
              "values": [
                "v1"
              ]
            }
          ]
      }
      """
    Then the HTTP status code should be "400"
    And I should receive the following response:
      """
      missing mandatory fields: entities
      """

  Scenario: [Test #4] POST /v1/policies with empty entities
    Given I am an admin user
    When I POST "/v1/policies"
      """
      {
          "entities": [],
          "role": "r1",
          "conditions": [
            {
              "attributes": [
                "a1"
              ],
              "operator": "and",
              "values": [
                "v1"
              ]
            }
          ]
      }
      """
    Then the HTTP status code should be "400"
    And I should receive the following response:
      """
      missing mandatory fields: entities
      """


  Scenario: [Test #5] POST /v1/policies without role
    Given I am an admin user
    When I POST "/v1/policies"
      """
      {
          "entities": [
            "e1",
            "e2"
          ],
          "conditions": [
            {
              "attributes": [
                "a1"
              ],
              "operator": "and",
              "values": [
                "v1"
              ]
            }
          ]
      }
      """
    Then the HTTP status code should be "400"
    And I should receive the following response:
      """
      missing mandatory fields: role
      """


  Scenario: [Test #6] POST /v1/policies with empty role
    Given I am an admin user
    When I POST "/v1/policies"
      """
      {
          "entities": ["e1"],
          "role": "",
          "conditions": [
            {
              "attributes": [
                "a1"
              ],
              "operator": "and",
              "values": [
                "v1"
              ]
            }
          ]
      }
      """
    Then the HTTP status code should be "400"
    And I should receive the following response:
      """
      missing mandatory fields: role
      """

  Scenario: [Test #7] POST /v1/policies without the correct permissions - the response status is 403 (forbidden)
    Given I am a viewer user
    When I POST "/v1/policies"
      """
      {
          "entities": ["e1"],
          "role": "",
          "conditions": [
            {
              "attributes": [
                "a1"
              ],
              "operator": "and",
              "values": [
                "v1"
              ]
            }
          ]
      }
      """
    Then the HTTP status code should be "403"
