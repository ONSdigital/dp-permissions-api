Feature: POST /v1/policies/{id} endpoint

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
                    "condition": {}
                }
            ]
            """

    Scenario: [Test #1] POST /v1/policies/{id} with all parameters returns 201
        Given I am an admin user
        When I POST "/v1/policies/new-policy"
            """
            {
                "entities": [
                    "e1",
                    "e2"
                ],
                "role": "r1",
                "condition": {
                    "attribute": "a1",
                    "operator": "StringEquals",
                    "values": [
                        "v1"
                    ]
                }
            }
            """
        Then the HTTP status code should be "201"
        And I should receive the following JSON response:
            """
            {
                "id": "new-policy",
                "entities": [
                    "e1",
                    "e2"
                ],
                "role": "r1",
                "condition": {
                    "attribute": "a1",
                    "operator": "StringEquals",
                    "values": [
                        "v1"
                    ]
                }
            }
            """
    
    Scenario: [Test #2] POST /v1/policies/{id} with missing entities returns 400
        Given I am an admin user
        When I POST "/v1/policies/new-policy"
            """
            {
                "role": "r1",
                "condition": {
                    "attribute": "a1",
                    "operator": "StringEquals",
                    "values": [
                        "v1"
                    ]
                }
            }
            """
        Then the HTTP status code should be "400"
        And I should receive the following JSON response:
            """
            {
                "errors": [
                    {
                        "code": "InvalidPolicyError",
                        "description": "missing mandatory fields: entities"
                    }
                ]
            }
            """
    
    Scenario: [Test #3] POST /v1/policies/{id} with an invalid Authorization header returns 401
        Given I am not authorised
        When I POST "/v1/policies/new-policy"
            """
            {
                "entities": [
                    "e1",
                    "e2"
                ],
                "role": "r1",
                "condition": {
                    "attribute": "a1",
                    "operator": "StringEquals",
                    "values": [
                        "v1"
                    ]
                }
            }
            """
        Then the HTTP status code should be "401"

    Scenario: [Test #4] POST /v1/policies/{id} without the correct permissions returns 403
        Given I am a viewer user
        When I POST "/v1/policies/new-policy"
            """
            {
                "entities": [
                    "e1",
                    "e2"
                ],
                "role": "r1",
                "condition": {
                    "attribute": "a1",
                    "operator": "StringEquals",
                    "values": [
                        "v1"
                    ]
                }
            }
            """
        Then the HTTP status code should be "403"

    Scenario: [Test #5] POST /v1/policies/{id} with an existing id returns 409
        Given I am an admin user
        When I POST "/v1/policies/admin"
            """
            {
                "entities": [
                    "e1",
                    "e2"
                ],
                "role": "r1",
                "condition": {
                    "attribute": "a1",
                    "operator": "StringEquals",
                    "values": [
                        "v1"
                    ]
                }
            }
            """
        Then the HTTP status code should be "409"
        And I should receive the following JSON response:
            """
            {
                "errors": [
                    {
                        "code": "PolicyAlreadyExistsError",
                        "description": "policy already exists with given ID"
                    }
                ]
            }
            """