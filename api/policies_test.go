package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1", "conditions": [{"attributes": ["a1"], "operator": "and", "values": ["v1"]}]}`)
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
					{Attributes: []string{"a1"}, Values: []string{"v1"}, Operator: "and"}},
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

func TestFailedAddPoliciesWithEmptyFields(t *testing.T) {
	t.Parallel()

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

		reader := strings.NewReader(`{"entities": ["e1", "e2"], "conditions": [{"attributes": ["a1"], "operator": "and", "values": ["v1"]}]}`)
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
			So(response, ShouldContainSubstring, "failed to parse json body")
		})
		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})
	})

}

func TestFailedAddPoliciesWhenPermissionStoreFails(t *testing.T) {
	Convey("when a permission store fails to insert a policy to data store", t, func() {

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
			So(response, ShouldContainSubstring, "Something went wrong")
		})

		Convey("Then the request body has been drained", func() {
			bytesRead, err := request.Body.Read(make([]byte, 1))
			So(bytesRead, ShouldEqual, 0)
			So(err, ShouldEqual, io.EOF)
		})

	})

}
