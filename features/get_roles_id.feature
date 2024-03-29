Feature: Behaviour of application when doing the GET /v1/roles/{id} endpoint, using a stripped down version of the database

  Scenario: [Test #1] GET /v1/roles/admin
    Given I have these roles:
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
    And I am an admin user
    When I GET "/v1/roles/admin"
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

  Scenario: [Test #2] Receive not found when doing a GET for a non existant role
    Given I have these roles:
      """
      [ ]
      """
    And I am an admin user
    When I GET "/v1/roles/unknown"
    Then the HTTP status code should be "404"

  Scenario: [Test #3] GET /v1/roles/admin with incorrect permissions - the response status is 403 (forbidden)
  Given I am a basic user
  When I GET "/v1/roles/admin"
  Then the HTTP status code should be "403"
