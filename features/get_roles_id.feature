Feature: Behaviour of application when doing the GET /roles/{id} endpoint, using a stripped down version of the database

  Scenario: [Test #3] GET /roles/admin
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
                  }
            ]
            """
    When I GET "/roles/admin"
    Then the HTTP status code should be "200"
    And the response header "Content-Type" should be "application/json; charset=utf-8"
    And I should receive the following JSON response:
            """
            {
              "id": "admin",
              "name": "Admin",
              "permissions": [
                "CreateRole",
                "Edit",
                "ReadOnly"
              ]
            }
            """

  Scenario: [Test #4] Receive not found when doing a GET for a non existant role
    Given I have this roles:
            """
            [ ]
            """
    When I GET "/roles/unknown"
    Then the HTTP status code should be "404"
    And the response header "Content-Type" should be "text/plain; charset=utf-8"
    And I should receive the following response:
            """
            role not found
            """"""
