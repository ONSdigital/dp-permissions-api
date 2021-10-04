Feature: Behaviour of application when doing the GET /v1/policies endpoint, using a stripped down version of the database

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
                            "operator": "=",
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

  Scenario: [Test #1] GET /v1/policies/viewer to fetch a policy having all parameters
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
                            "operator": "=",
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

  Scenario: [Test #2] GET /v1/policies/admin to fetch a policy without conditions
    When I GET "/v1/policies/admin"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
            """
            {
                    "id": "admin",
                    "role": "admin",
                    "entities": [
                      "group/admin"
                    ]
                }
            """

  Scenario: [Test #3] Receive not found when doing a GET for a non existent policy
    When I GET "/v1/policies/notFound"
    Then the HTTP status code should be "404"
    And the response header "Content-Type" should be "text/plain; charset=utf-8"
    And I should receive the following response:
            """
            policy not found
            """"""
