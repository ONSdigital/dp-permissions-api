package sdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	dphttp "github.com/ONSdigital/dp-net/http"
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
					Body:       ioutil.NopCloser(bytes.NewReader(getExampleBundleJson())),
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
				So(policy.Condition.Operator, ShouldEqual, "StringEquals")
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
						Body:       ioutil.NopCloser(strings.NewReader(`bad response`)),
					}, nil
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(getExampleBundleJson())),
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
				So(policy.Condition.Operator, ShouldEqual, "StringEquals")
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
					Body:       ioutil.NopCloser(strings.NewReader(`bad response`)),
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

func getExampleBundleJson() []byte {
	bundle := getExampleBundle()
	permissionsBundleJson, err := json.Marshal(bundle)
	So(err, ShouldBeNil)
	return permissionsBundleJson
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
