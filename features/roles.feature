Feature: Behaviour of application when doing the GET /roles endpoint, using a stripped down version of the database

    # A Background applies to all scenarios in this Feature
  Background:
    Given I have this roles:
            """
            [
                {
                    "id": "admin",
                    "name": "Admin",
                    "permissions": [
                      "CreateRole",
                      "Edit",
                      "ReadOnly"
                    ]
                  },
                  {
                    "id": "publisher",
                    "name": "Publisher",
                    "permissions": [
                      "Edit",
                      "ReadOnly"
                    ]
                  },
                  {
                    "id": "readonly",
                    "name": "Readonly",
                    "permissions": [
                        "ReadOnly"
                    ]
                  }
            ]
            """

  Scenario: [Test #1] GET /roles with default offset and limit
    When I GET "/roles"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
            """
            {
                "count": 3,
                "offset": 0,
                "limit": 20,
                "items": [
                  {
                    "id": "admin",
                    "name": "Admin",
                    "permissions": [
                      "CreateRole",
                      "Edit",
                      "ReadOnly"
                    ]
                  },
                  {
                    "id": "publisher",
                    "name": "Publisher",
                    "permissions": [
                      "Edit",
                      "ReadOnly"
                    ]
                  },
                  {
                    "id": "readonly",
                    "name": "Readonly",
                    "permissions": [
                        "ReadOnly"
                    ]
                  }
                ],
                "total_count": 3
            }
            """

  Scenario: [Test #2] GET /roles with offset and limit
    When I GET "/roles?offset=2&limit=1"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
            """
            {
              "count": 1,
              "offset": 2,
              "limit": 1,
              "items": [
                {
                  "id": "readonly",
                  "name": "Readonly",
                  "permissions": [
                    "ReadOnly"
                  ]
                }
              ],
              "total_count": 3
            }
            """"""
