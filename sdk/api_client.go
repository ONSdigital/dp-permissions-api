package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
)

// package level constants
const (
	bundlerEndpoint = "%s/v1/permissions-bundle"
)

// HTTPClient is the interface that defines a client for making HTTP requests
type HTTPClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// APIClient implementation of permissions.Store that gets permission data from the permissions API
type APIClient struct {
	host            string
	httpCli         HTTPClient
	backoffSchedule []time.Duration
	options         Options
}

// Options is a struct containing for customised options for the API client
type Options struct {
}

// NewClient constructs a new APIClient instance with a default http client and Options.
func NewClient(host string) *APIClient {
	return NewClientWithOptions(host, Options{})
}

// NewClientWithOptions returns a new APIClient with default http
func NewClientWithOptions(host string, opts Options) *APIClient {
	return NewClientWithClienter(host, dphttp.NewClient(), opts)
}

// NewClientWithClienter constructs a new APIClient instance.
func NewClientWithClienter(host string, httpClient HTTPClient, opts Options) *APIClient {
	return &APIClient{
		host:    host,
		httpCli: httpClient,
		options: opts,
	}
}

// GetPermissionsBundle gets the permissions bundle data from the permissions API.
func (c *APIClient) GetPermissionsBundle(ctx context.Context) (Bundle, error) {

	uri := fmt.Sprintf(bundlerEndpoint, c.host)

	log.Info(ctx, "GetPermissionsBundle: starting permissions bundle request", log.Data{"uri": uri})

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Info(ctx, "GetPermissionsBundle: error building new request", log.Data{"err": err.Error()})
		return nil, err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "GetPermissionsBundle: error executing request", log.Data{"err": err.Error()})
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	log.Info(ctx, "GetPermissionsBundle: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status returned from the permissions api permissions-bundle endpoint: %s", resp.Status)
	}

	permissions, err := getPermissionsBundleFromResponse(resp.Body)
	if err != nil {
		log.Info(ctx, "GetPermissionsBundle: error getting permissions bundle from response", log.Data{"err": err.Error()})
		return nil, err
	}

	log.Info(ctx, "GetPermissionsBundle: returning requested permissions to caller")

	return permissions, nil
}

func getPermissionsBundleFromResponse(reader io.Reader) (Bundle, error) {
	b, err := getResponseBytes(reader)
	if err != nil {
		return nil, err
	}

	var bundle Bundle

	if err := json.Unmarshal(b, &bundle); err != nil {
		return nil, ErrFailedToParsePermissionsResponse
	}

	return bundle, nil
}

func getResponseBytes(reader io.Reader) ([]byte, error) {
	if reader == nil {
		return nil, ErrGetPermissionsResponseBodyNil
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if b == nil || len(b) == 0 {
		return nil, ErrGetPermissionsResponseBodyNil
	}

	return b, nil
}
