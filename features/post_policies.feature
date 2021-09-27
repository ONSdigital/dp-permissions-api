Feature: Behaviour of application when doing the POST /v1/policies endpoint, using a stripped down version of the database

  Scenario: [Test #1] POST /v1/policies with all the parameters
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


  Scenario: [Test #1] POST /v1/policies without conditions
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

  Scenario: [Test #2] POST /v1/policies without entities
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
  Scenario: [Test #2] POST /v1/policies with empty entities
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


  Scenario: [Test #2] POST /v1/policies without role
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


  Scenario: [Test #2] POST /v1/policies with empty role
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
