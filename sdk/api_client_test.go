package sdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/models"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-permissions-api/sdk"
	. "github.com/smartystreets/goconvey/convey"
)

var host = "localhost:1234"

func TestNewClient(t *testing.T) {
	Convey("Given some host", t, func() {
		someHost := host
		Convey("When NewClient is called", func() {
			apiClient := sdk.NewClient(someHost)
			Convey("Then the api is not nil", func() {
				So(apiClient, ShouldNotBeNil)
			})
		})
	})
}

func TestNewClientWithOptions(t *testing.T) {
	Convey("Given some host and custom options", t, func() {
		someHost := host
		ops := sdk.Options{}
		Convey("When NewClientWithOptions is called", func() {
			apiClient := sdk.NewClientWithOptions(someHost, ops)
			Convey("Then the api is not nil", func() {
				So(apiClient, ShouldNotBeNil)
			})
		})
	})
}

func TestAPIClient_GetPermissionsBundle(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a successful permissions bundle response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(getExampleBundleJSON())),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPermissionsBundle is called", func() {
			bundle, err := apiClient.GetPermissionsBundle(ctx)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected permissions bundle is returned", func() {
				So(bundle, ShouldNotBeNil)

				policies := bundle["permission/admin"]["group/admin"]
				So(policies, ShouldHaveLength, 1)

				policy := policies[0]
				So(policy.ID, ShouldEqual, "policy/123")
				So(policy.Condition.Attribute, ShouldEqual, "collection_id")
				So(policy.Condition.Operator, ShouldEqual, sdk.Operator("StringEquals"))
				So(policy.Condition.Values, ShouldHaveLength, 1)
				So(policy.Condition.Values[0], ShouldEqual, "col123")
			})
		})
	})
}

func TestAPIClient_GetPermissionsBundle_SucceedsOnSecondAttempt(t *testing.T) {
	ctx := context.Background()

	// GetPermissionsBundle request counter
	retryCount := 1

	Convey("Given a mock http client that returns a successful permissions bundle response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				if retryCount == 1 {
					retryCount++
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(strings.NewReader(`bad response`)),
					}, nil
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(getExampleBundleJSON())),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPermissionsBundle is called", func() {
			bundle, err := apiClient.GetPermissionsBundle(ctx)

			Convey("Bundle request made twice (failed on first attempt)", func() {
				So(retryCount, ShouldEqual, 2)
			})

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected permissions bundle is returned", func() {
				So(bundle, ShouldNotBeNil)

				policies := bundle["permission/admin"]["group/admin"]
				So(policies, ShouldHaveLength, 1)

				policy := policies[0]
				So(policy.ID, ShouldEqual, "policy/123")
				So(policy.Condition.Attribute, ShouldEqual, "collection_id")
				So(policy.Condition.Operator, ShouldEqual, sdk.Operator("StringEquals"))
				So(policy.Condition.Values, ShouldHaveLength, 1)
				So(policy.Condition.Values[0], ShouldEqual, "col123")
			})
		})
	})
}

func TestAPIClient_GetPermissionsBundle_HTTPError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("something went wrong")

	Convey("Given a mock http client that returns an error", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, expectedErr
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPermissionsBundle is called", func() {
			bundle, err := apiClient.GetPermissionsBundle(ctx)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldEqual, expectedErr)
			})

			Convey("Then the permissions bundle is nil", func() {
				So(bundle, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_GetPermissionsBundle_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code other than 200", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 internal server error",
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPermissionsBundle is called", func() {
			bundle, err := apiClient.GetPermissionsBundle(ctx)

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, "unexpected status returned from the permissions api permissions-bundle endpoint: 500 internal server error")
			})

			Convey("Then the permissions bundle is nil", func() {
				So(bundle, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_GetPermissionsBundle_NilResponseBody(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response with a nil body", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPermissionsBundle is called", func() {
			bundle, err := apiClient.GetPermissionsBundle(ctx)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldEqual, sdk.ErrGetPermissionsResponseBodyNil)
			})

			Convey("Then the permissions bundle is nil", func() {
				So(bundle, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_GetPermissionsBundle_UnexpectedResponseBody(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response with unexpected body content", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`bad response`)),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPermissionsBundle is called", func() {
			bundle, err := apiClient.GetPermissionsBundle(ctx)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldEqual, sdk.ErrFailedToParsePermissionsResponse)
			})

			Convey("Then the permissions bundle is nil", func() {
				So(bundle, ShouldBeNil)
			})
		})
	})
}

func getExampleBundleJSON() []byte {
	bundle := getExampleBundle()
	permissionsBundleJSON, err := json.Marshal(bundle)
	So(err, ShouldBeNil)
	return permissionsBundleJSON
}

func getExampleBundle() sdk.Bundle {
	bundle := sdk.Bundle{
		"permission/admin": map[string][]sdk.Policy{
			"group/admin": {
				{
					ID: "policy/123",
					Condition: sdk.Condition{
						Attribute: "collection_id",
						Operator:  "StringEquals",
						Values:    []string{"col123"}},
				},
			},
		},
	}
	return bundle
}

// == Roles ===

func TestAPIClient_GetRoles(t *testing.T) {
	ctx := context.Background()
	result := models.Roles{
		Count:  2,
		Offset: 0,
		Limit:  2,
		Items: []models.Role{{
			ID:          "1",
			Name:        "test",
			Permissions: []string{"all"},
		}, {
			ID:          "2",
			Name:        "test",
			Permissions: []string{"all"},
		}},
		TotalCount: 0,
	}

	bresult, err := json.Marshal(result)
	if err != nil {
		t.Failed()
	}

	Convey("Given a mock http client that returns a successful all roles response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bresult)),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRoles is called", func() {
			roles, err := apiClient.GetRoles(ctx)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected roles is returned", func() {
				So(roles.Items[0], ShouldResemble, models.Role{ID: "1", Name: "test", Permissions: []string{"all"}})
				So(roles.Items[1], ShouldResemble, models.Role{ID: "2", Name: "test", Permissions: []string{"all"}})
			})
		})
	})
}

func TestAPIClient_GetRoles_BadRequest(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a error in response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, errors.New("bad request")
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRoles is called", func() {
			_, err := apiClient.GetRoles(ctx)

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("bad request"))
			})
		})
	})
}

func TestAPIClient_GetRoles_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code 400", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status: `Invalid request, reasons can be one of the following:
              * query parameters incorrect offset provided
              * query parameters incorrect limit provided`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRoles is called", func() {
			roles, err := apiClient.GetRoles(ctx)

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getallpolicies endpoint: Invalid request, reasons can be one of the following:
              * query parameters incorrect offset provided
              * query parameters incorrect limit provided`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(roles, ShouldBeNil)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 403", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status: `Unauthorised request, reason is:
              * Requestor does not have necessary permissions to access this resource`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRoles is called", func() {
			roles, err := apiClient.GetRoles(ctx)

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getallpolicies endpoint: Unauthorised request, reason is:
              * Requestor does not have necessary permissions to access this resource`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(roles, ShouldBeNil)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 500", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     `Failed to process the request due to an internal error`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRoles is called", func() {
			roles, err := apiClient.GetRoles(ctx)

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getallpolicies endpoint: Failed to process the request due to an internal error`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(roles, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_GetRole(t *testing.T) {
	ctx := context.Background()
	result := models.Roles{
		Count:  2,
		Offset: 0,
		Limit:  2,
		Items: []models.Role{{
			ID:          "2",
			Name:        "test",
			Permissions: []string{"all"},
		}},
		TotalCount: 0,
	}

	bresult, err := json.Marshal(result)
	if err != nil {
		t.Failed()
	}

	Convey("Given a mock http client that returns a successful role response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bresult)),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			role, err := apiClient.GetRole(ctx, "2")

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected role is returned", func() {
				So(role.Items[0], ShouldResemble, models.Role{ID: "2", Name: "test", Permissions: []string{"all"}})
			})
		})
	})
}

func TestAPIClient_GetRole_BadRequest(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a error in response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, errors.New("bad request")
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			_, err := apiClient.GetRole(ctx, "1")

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("bad request"))
			})
		})
	})
}

func TestAPIClient_GetRole_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code 400", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     `Invalid request`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			role, err := apiClient.GetRole(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getrole endpoint: Invalid request`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(role, ShouldBeNil)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 403", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     `Unauthorised request`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			role, err := apiClient.GetRole(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getrole endpoint: Unauthorised request`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(role, ShouldBeNil)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 404", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     `Requested id can not be found`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			role, err := apiClient.GetRole(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getrole endpoint: Requested id can not be found`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(role, ShouldBeNil)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 500", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     `Failed to process the request due to an internal error`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			role, err := apiClient.GetRole(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getrole endpoint: Failed to process the request due to an internal error`)
			})

			Convey("Then the permissions roles is nil", func() {
				So(role, ShouldBeNil)
			})
		})
	})
}

// == Policy ===

func TestAPIClient_PostPolicy(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a successful add policy response", t, func() {
		result := models.Policy{
			ID:        "",
			Entities:  nil,
			Role:      "",
			Condition: models.Condition{},
		}

		bresult, err := json.Marshal(result)
		if err != nil {
			t.Failed()
		}

		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bresult)),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PostPolicy is called", func() {
			policy, err := apiClient.PostPolicy(ctx, models.PolicyInfo{
				Entities:  nil,
				Role:      "1",
				Condition: models.Condition{},
			})

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Policy should not be empty", func() {
				So(policy, ShouldResemble, &result)
			})
		})
	})
}

func TestAPIClient_PostPolicy_BadRequest(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a error in response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, errors.New("bad request")
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PostPolicy is called", func() {
			_, err := apiClient.PostPolicy(ctx, models.PolicyInfo{})

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("bad request"))
			})
		})
	})
}

func TestAPIClient_PostPolicy_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code 400", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     `Bad request. Invalid policy supplied`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			_, err := apiClient.PostPolicy(ctx, models.PolicyInfo{})

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-addpolicy endpoint: Bad request. Invalid policy supplied`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 403", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     `Unauthorised request`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetRole is called", func() {
			_, err := apiClient.PostPolicy(ctx, models.PolicyInfo{})

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-addpolicy endpoint: Unauthorised request`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 500", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     `Failed to process the request due to an internal error`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetAllRoles is called", func() {
			_, err := apiClient.PostPolicy(ctx, models.PolicyInfo{})

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-addpolicy endpoint: Failed to process the request due to an internal error`)
			})
		})
	})
}

func TestAPIClient_DeletePolicy(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a successful delete policy response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When DeletePolicy is called", func() {
			err := apiClient.DeletePolicy(ctx, "1")

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_DeletePolicy_BadRequest(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a error in response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, errors.New("bad request")
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When DeletePolicy is called", func() {
			err := apiClient.DeletePolicy(ctx, "1")

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("bad request"))
			})
		})
	})
}

func TestAPIClient_DeletePolicy_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code 400", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     `Bad request. Invalid policy supplied`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			err := apiClient.DeletePolicy(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-deletepolicy endpoint: Bad request. Invalid policy supplied`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 403", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     `Unauthorised request`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			err := apiClient.DeletePolicy(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-deletepolicy endpoint: Unauthorised request`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 500", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     `Failed to process the request due to an internal error`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			err := apiClient.DeletePolicy(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-deletepolicy endpoint: Failed to process the request due to an internal error`)
			})
		})
	})
}

func TestAPIClient_GetPolicy(t *testing.T) {
	ctx := context.Background()

	result := models.Policy{
		ID:        "1",
		Entities:  nil,
		Role:      "1",
		Condition: models.Condition{},
	}

	bresult, err := json.Marshal(result)
	if err != nil {
		t.Failed()
	}

	Convey("Given a mock http client that returns a successful get policy response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bresult)),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			policy, err := apiClient.GetPolicy(ctx, "1")

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected policy is returned", func() {
				So(policy, ShouldResemble, &result)
			})
		})
	})
}

func TestAPIClient_GetPolicy_BadRequest(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a error in response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, errors.New("bad request")
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			policy, err := apiClient.GetPolicy(ctx, "1")

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("bad request"))
			})

			Convey("Then policy should be nil", func() {
				So(policy, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_GetPolicy_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code 400", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     `Bad request. Invalid policy supplied`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			_, err := apiClient.GetPolicy(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getpolicy endpoint: Bad request. Invalid policy supplied`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 403", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     `Unauthorised request`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			_, err := apiClient.GetPolicy(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getpolicy endpoint: Unauthorised request`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 500", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     `Failed to process the request due to an internal error`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When GetPolicy is called", func() {
			_, err := apiClient.GetPolicy(ctx, "1")

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-getpolicy endpoint: Failed to process the request due to an internal error`)
			})
		})
	})
}

func TestAPIClient_PutPolicy(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a successful put policy response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PutPolicy is called", func() {
			err := apiClient.PutPolicy(ctx, "1", models.Policy{
				ID:        "",
				Entities:  nil,
				Role:      "",
				Condition: models.Condition{},
			})

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestAPIClient_PutPolicy_BadRequest(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a error in response", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return nil, errors.New("bad request")
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PutPolicy is called", func() {
			err := apiClient.PutPolicy(ctx, "1", models.Policy{})

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("bad request"))
			})
		})
	})
}

func TestAPIClient_PutPolicy_Non200ResponseCodeReturned(t *testing.T) {
	ctx := context.Background()

	Convey("Given a mock http client that returns a response code 400", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     `Bad request. Invalid policy supplied`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PutPolicy is called", func() {
			err := apiClient.PutPolicy(ctx, "", models.Policy{})

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-putpolicy endpoint: Bad request. Invalid policy supplied`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 403", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     `Unauthorised request`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PutPolicy is called", func() {
			err := apiClient.PutPolicy(ctx, "", models.Policy{})

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-putpolicy endpoint: Unauthorised request`)
			})
		})
	})

	Convey("Given a mock http client that returns a response code 500", t, func() {
		httpClient := &dphttp.ClienterMock{
			DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     `Failed to process the request due to an internal error`,
				}, nil
			},
		}
		apiClient := sdk.NewClientWithClienter(host, httpClient, sdk.Options{})

		Convey("When PutPolicy is called", func() {
			err := apiClient.PutPolicy(ctx, "", models.Policy{})

			Convey("Then the expected error is returned", func() {
				So(err.Error(), ShouldEqual, `unexpected status returned from the permissions api permissions-putpolicy endpoint: Failed to process the request due to an internal error`)
			})
		})
	})
}
