Feature: GET /v1/permissions-bundle endpoint

  Background:
    Given I have these roles:
            """
            [
                {
                    "id": "admin",
                    "name": "Admin",
                    "permissions": [
                      "legacy.read", "legacy.update", "users.add"
                    ]
                  },
                  {
                    "id": "publisher",
                    "name": "Publisher",
                    "permissions": [
                      "legacy.read", "legacy.update"
                    ]
                  },
                  {
                    "id": "viewer",
                    "name": "Viewer",
                    "permissions": [
                        "legacy.read"
                    ]
                  }
            ]
            """
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
                },
                {
                    "id": "publisher",
                    "role": "publisher",
                    "entities": [
                      "group/publisher"
                    ],
                    "condition": {}
                },
                {
                    "id": "viewer",
                    "role": "viewer",
                    "entities": [
                      "group/viewer"
                    ],
                    "condition": {
                            "operator": "StringEquals",
                            "attribute": "collection-id",
                            "values": [
                              "collection-765"
                            ]
                    }
                }
            ]
            """

  Scenario: GET /v1/permissions-bundle
    When I GET "/v1/permissions-bundle"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
            """
            {
              "legacy.read": {
                "group/admin": [
                  {
                    "id": "admin"
                  }
                ],
                "group/publisher": [
                  {
                    "id": "publisher"
                  }
                ],
                "group/viewer": [
                  {
                    "id": "viewer",
                    "condition":{
                        "attribute": "collection-id",
                        "operator": "StringEquals",
                        "values": [
                          "collection-765"
                        ]
                    }
                  }
                ]
              },
              "legacy.update": {
                "group/admin": [
                  {
                    "id": "admin"
                  }
                ],
                "group/publisher": [
                  {
                    "id": "publisher"
                  }
                ]
              },
              "users.add": {
                "group/admin": [
                  {
                    "id": "admin"
                  }
                ]
              }
            }
            """
