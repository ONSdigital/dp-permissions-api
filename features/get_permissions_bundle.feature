Feature: GET /v1/permissions-bundle endpoint

  Background:
    Given I have this roles:
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
                    "id": "admin",
                    "entities": [
                      "group/admin"
                    ],
                    "role": "admin"
                  }
                ],
                "group/publisher": [
                  {
                    "id": "publisher",
                    "entities": [
                      "group/publisher"
                    ],
                    "role": "publisher"
                  }
                ],
                "group/viewer": [
                  {
                    "id": "viewer",
                    "entities": [
                      "group/viewer"
                    ],
                    "role": "viewer",
                    "conditions": [
                      {
                        "attributes": [
                          "collection-id"
                        ],
                        "operator": "=",
                        "values": [
                          "collection-765"
                        ]
                      }
                    ]
                  }
                ]
              },
              "legacy.update": {
                "group/admin": [
                  {
                    "id": "admin",
                    "entities": [
                      "group/admin"
                    ],
                    "role": "admin"
                  }
                ],
                "group/publisher": [
                  {
                    "id": "publisher",
                    "entities": [
                      "group/publisher"
                    ],
                    "role": "publisher"
                  }
                ]
              },
              "users.add": {
                "group/admin": [
                  {
                    "id": "admin",
                    "entities": [
                      "group/admin"
                    ],
                    "role": "admin"
                  }
                ]
              }
            }
            """
