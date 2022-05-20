package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testPolicyID = "testPoliciesID"
)

func TestSuccessfulAddPolicies(t *testing.T) {
	t.Parallel()

	Convey("Given a permissions store", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			AddPolicyFunc: func(ctx context.Context, policy *models.Policy) (*models.Policy, error) {
				if policy.Entities != nil {
					policy.ID = testPolicyID
					return policy, nil
				}
				return nil, errors.New("Something went wrong")
			},
		}

		permissionsApi := setupAPIWithStore(mockedPermissionsStore)

		Convey("When a POST request is made to the policies endpoint with all the policies properties", func() {
			reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1", "conditions": [{"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}]}`)
			request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)

			Convey("Then the permissions store is called to create a new policy", func() {
				So(len(mockedPermissionsStore.AddPolicyCalls()), ShouldEqual, 1)
			})

			Convey("Then the response is 201 created", func() {
				So(responseWriter.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("Then the request body has been drained", func() {
				bytesRead, err := request.Body.Read(make([]byte, 1))
				So(bytesRead, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})

			Convey("Then the response body has newly created policy", func() {
				policy := models.Policy{}
				json.Unmarshal(responseWriter.Body.Bytes(), &policy)

				So(policy, ShouldNotBeNil)
				So(policy.ID, ShouldEqual, testPolicyID)
				So(policy.Role, ShouldResemble, "r1")
				So(policy.Entities, ShouldResemble, []string{"e1", "e2"})
				So(policy.Conditions, ShouldResemble, []models.Condition{
					{Attribute: "a1", Values: []string{"v1"}, Operator: models.OperatorStringEquals}},
				)
			})
		})

		Convey("When a POST request is made to the policies endpoint without conditions", func() {
			reader := strings.NewReader(`{"entities": ["e1"], "role": "r1"}`)
			request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)

			Convey("Then the permissions store is called to create a new policy", func() {
				So(len(mockedPermissionsStore.AddPolicyCalls()), ShouldEqual, 1)
			})

			Convey("Then the response is 201 created", func() {
				So(responseWriter.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("Then the request body has been drained", func() {
				bytesRead, err := request.Body.Read(make([]byte, 1))
				So(bytesRead, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})

			Convey("Then the response body has newly created policy", func() {
				policy := models.Policy{}
				json.Unmarshal(responseWriter.Body.Bytes(), &policy)

				So(policy, ShouldNotBeNil)
				So(policy.ID, ShouldEqual, testPolicyID)
				So(policy.Role, ShouldResemble, "r1")
				So(policy.Entities, ShouldResemble, []string{"e1"})
				So(policy.Conditions, ShouldResemble, []models.Condition(nil))
			})
		})
	})

}

func TestFailedAddPoliciesWithInvalidPolicy(t *testing.T) {
	t.Parallel()

	Convey("When a POST request is made to the policies endpoint without entities", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{"role": "r1"}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "missing mandatory fields: entities")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

	Convey("When a POST request is made to the policies endpoint with empty entities", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{"entities": [], "role": "r1"}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "missing mandatory fields: entities")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

	Convey("When a POST request is made to the policies without a role", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{"entities": ["e1", "e2"], "conditions": [{"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}]}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "missing mandatory fields: role")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

	Convey("When a POST request is made to the policies with empty role", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "", "conditions": [{"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}]}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "missing mandatory fields: role")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

	Convey("When a POST request is made to the policies with an invalid condition operator", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1", "conditions": [{"attribute": "a1", "operator": "And", "values": ["v1"]}]}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "invalid field values: condition operator And")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})
}

func TestFailedAddPoliciesWithBadJson(t *testing.T) {
	t.Parallel()

	Convey("When a POST request is made to the policies endpoint with an empty JSON message", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "missing mandatory fields: entities, role")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

	Convey("When a POST request is made to the policies endpoint with an invalid JSON message", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, models.UnmarshalFailedDescription)
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})
}

func TestFailedAddPoliciesWhenPermissionStoreFails(t *testing.T) {
	Convey("When a permission store fails to insert a policy to data store", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			AddPolicyFunc: func(ctx context.Context, policy *models.Policy) (*models.Policy, error) {
				return nil, errors.New("Something went wrong")
			},
		}

		permissionsApi := setupAPIWithStore(mockedPermissionsStore)

		reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1"}`)
		request, _ := http.NewRequest("POST", "http://localhost:25400/v1/policies", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the permissions store is called to create a new policy", func() {
			So(len(mockedPermissionsStore.AddPolicyCalls()), ShouldEqual, 1)
		})

		Convey("Then the response is 500 internal server error, , with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusInternalServerError)

			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, models.InternalServerErrorDescription)
		})

		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})

	})
}

func TestGetPolicyHandler(t *testing.T) {
	Convey("Given a GetPolicy Handler", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			GetPolicyFunc: func(ctx context.Context, id string) (*models.Policy, error) {
				switch id {
				case testPolicyID:
					return &models.Policy{
						ID:         testPolicyID,
						Entities:   []string{"e1", "e2"},
						Role:       "r1",
						Conditions: []models.Condition{{Attribute: "al", Operator: models.OperatorStringEquals, Values: []string{"v1"}}}}, nil
				case "NOTFOUND":
					return nil, apierrors.ErrPolicyNotFound
				default:
					return nil, errors.New("Something went wrong")
				}
			},
		}

		permissionsApi := setupAPIWithStore(mockedPermissionsStore)

		Convey("When an existing policy is requested with its policy ID", func() {

			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25400/v1/policies/%s", testPolicyID), nil)
			responseRecorder := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseRecorder, request)

			Convey("The matched policy is returned with status code 200", func() {
				expectedPolicy := models.Policy{
					ID:         testPolicyID,
					Entities:   []string{"e1", "e2"},
					Role:       "r1",
					Conditions: []models.Condition{{Attribute: "al", Operator: models.OperatorStringEquals, Values: []string{"v1"}}}}

				policy := models.Policy{}
				payload, _ := ioutil.ReadAll(responseRecorder.Body)
				err := json.Unmarshal(payload, &policy)
				So(err, ShouldBeNil)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
				So(policy, ShouldResemble, expectedPolicy)
			})
		})

		Convey("When a non existing policy id is requested a Not Found response with 404 status code is returned", func() {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/policies/NOTFOUND", nil)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)
			response := responseWriter.Body.String()

			So(responseWriter.Code, ShouldEqual, http.StatusNotFound)
			So(response, ShouldContainSubstring, models.PolicyNotFoundDescription)
		})

		Convey("When a failed to fetch the policy from DB should return a status code of 500", func() {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/policies/XYZ", nil)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)
			response := responseWriter.Body.String()

			So(responseWriter.Code, ShouldEqual, http.StatusInternalServerError)
			So(response, ShouldContainSubstring, models.InternalServerErrorDescription)
		})
	})
}

func TestSuccessfulUpdatePolicy(t *testing.T) {
	t.Parallel()

	Convey("Given a permissions store", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			UpdatePolicyFunc: func(ctx context.Context, policy *models.Policy) (*models.UpdateResult, error) {
				if policy.ID == "existing_policy" {
					return &models.UpdateResult{ModifiedCount: 1}, nil
				} else if policy.ID == "new_policy" {
					return &models.UpdateResult{UpsertedCount: 1}, nil
				}
				return nil, errors.New("Something went wrong")
			},
		}

		permissionsApi := setupAPIWithStore(mockedPermissionsStore)

		Convey("When a PUT request is made to the update policies endpoint to update an existing policy", func() {
			reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1", "conditions": [{"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}]}`)
			request, _ := http.NewRequest("PUT", "http://localhost:25400/v1/policies/existing_policy", reader)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)

			Convey("Then the permissions store is called to create a update policy", func() {
				So(len(mockedPermissionsStore.UpdatePolicyCalls()), ShouldEqual, 1)
			})

			Convey("Then the response is 200", func() {
				So(responseWriter.Code, ShouldEqual, http.StatusOK)
			})

			Convey("Then the request body has been drained", func() {
				bytesRead, err := request.Body.Read(make([]byte, 1))
				So(bytesRead, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})
		})

		Convey("When a PUT request is made to the update policies endpoint with a non-existing policy id", func() {
			reader := strings.NewReader(`{"entities": ["e1"], "role": "r1"}`)
			request, _ := http.NewRequest("PUT", "http://localhost:25400/v1/policies/new_policy", reader)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)

			Convey("Then the permissions store is called to upsert a policy", func() {
				So(len(mockedPermissionsStore.UpdatePolicyCalls()), ShouldEqual, 1)
			})

			Convey("Then the response is 201 created", func() {
				So(responseWriter.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("Then the request body has been drained", func() {
				bytesRead, err := request.Body.Read(make([]byte, 1))
				So(bytesRead, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}

func TestFailedUpdatePoliciesWithBadJson(t *testing.T) {
	t.Parallel()

	Convey("When a PUT request is made to the update policies endpoint with an empty JSON message", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{}`)
		request, _ := http.NewRequest("PUT", "http://localhost:25400/v1/policies/policyid", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, "missing mandatory fields: entities, role")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

	Convey("When a PUT request is made to the update policies endpoint with an invalid JSON message", t, func() {
		permissionsApi := setupAPI()

		reader := strings.NewReader(`{`)
		request, _ := http.NewRequest("PUT", "http://localhost:25400/v1/policies/policyid", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the response is 400 bad request, with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusBadRequest)
			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, models.UnmarshalFailedDescription)
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})
}

func TestFailedUpdatePoliciesWhenPermissionStoreFails(t *testing.T) {
	Convey("When a permission store fails to insert a policy to data store", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			UpdatePolicyFunc: func(ctx context.Context, policy *models.Policy) (*models.UpdateResult, error) {
				return nil, errors.New("Something went wrong")
			},
		}

		permissionsApi := setupAPIWithStore(mockedPermissionsStore)

		reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1"}`)
		request, _ := http.NewRequest("PUT", "http://localhost:25400/v1/policies/policyid", reader)
		responseWriter := httptest.NewRecorder()
		permissionsApi.Router.ServeHTTP(responseWriter, request)

		Convey("Then the permissions store is called to update a policy", func() {
			So(len(mockedPermissionsStore.UpdatePolicyCalls()), ShouldEqual, 1)
		})

		Convey("Then the response is 500 internal server error with the expected response body", func() {
			So(responseWriter.Code, ShouldEqual, http.StatusInternalServerError)

			response := responseWriter.Body.String()
			So(response, ShouldContainSubstring, models.InternalServerErrorDescription)
		})

		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})
}

func TestDeletePolicyHandler(t *testing.T) {
	Convey("Given a DeletePolicy Handler", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			DeletePolicyFunc: func(ctx context.Context, id string) error {
				switch id {
				case testPolicyID:
					return nil
				case "NOTFOUND":
					return apierrors.ErrPolicyNotFound
				default:
					return errors.New("Something went wrong")
				}
			},
			GetPolicyFunc: func(ctx context.Context, id string) (*models.Policy, error) {
				switch id {
				case testPolicyID:
					return &models.Policy{
						ID:         testPolicyID,
						Entities:   []string{"e1", "e2"},
						Role:       "r1",
						Conditions: []models.Condition{{Attribute: "al", Operator: models.OperatorStringEquals, Values: []string{"v1"}}}}, nil
				case "NOTFOUND":
					return nil, apierrors.ErrPolicyNotFound
				default:
					return nil, errors.New("Something went wrong")
				}
			},
		}

		permissionsApi := setupAPIWithStore(mockedPermissionsStore)

		Convey("When a DELETE request is made to an existing policy with its policy ID", func() {

			request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:25400/v1/policies/%s", testPolicyID), nil)
			responseRecorder := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseRecorder, request)

			Convey("Then the permissions store is called to delete a policy", func() {
				So(len(mockedPermissionsStore.DeletePolicyCalls()), ShouldEqual, 1)
			})

			Convey("The matched policy is returned with status code 204", func() {
				So(responseRecorder.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("When a DELETE request is made to a non existing policy id, a Not Found response with 404 status code is returned", func() {
			request := httptest.NewRequest(http.MethodDelete, "http://localhost:25400/v1/policies/NOTFOUND", nil)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)
			response := responseWriter.Body.String()

			So(responseWriter.Code, ShouldEqual, http.StatusNotFound)
			So(response, ShouldContainSubstring, models.PolicyNotFoundDescription)
		})

		Convey("When a failed DELETE request to the policy from the DB should return a status code of 500", func() {
			request := httptest.NewRequest(http.MethodDelete, "http://localhost:25400/v1/policies/XYZ", nil)
			responseWriter := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(responseWriter, request)
			response := responseWriter.Body.String()

			So(responseWriter.Code, ShouldEqual, http.StatusInternalServerError)
			So(response, ShouldContainSubstring, models.InternalServerErrorDescription)
		})
	})
}
