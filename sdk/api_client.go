package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

// package level constants
const (
	bundlerEndpoint   = "%s/v1/permissions-bundle"
	addPolicyEndpoint = "%s/v1/policies"    // Add policy
	policyEndpoint    = "%s/v1/policies/%s" // Delete policy
	rolesEndpoint     = "%s/v1/roles"       // Add roles
	getRoleEndpoint   = "%s/v1/roles/%s"    // Get roles
)

// HTTPClient is the interface that defines a client for making HTTP requests
type HTTPClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// APIClient implementation of permissions.Store that gets permission data from the permissions API
type APIClient struct {
	host    string
	httpCli HTTPClient
	options Options
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

// == Roles Endpoint ==

func (c *APIClient) GetRoles(ctx context.Context) (*models.Roles, error) {
	uri := fmt.Sprintf(rolesEndpoint, c.host)

	log.Info(ctx, "GetRoles: starting permissions get policy request", log.Data{"uri": uri})

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Info(ctx, "GetRoles: error building new request", log.Data{"err": err.Error()})
		return nil, err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "GetRoles: error executing request", log.Data{"err": err.Error()})
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status returned from the permissions api permissions-getallpolicies endpoint: %s", resp.Status)
	}

	log.Info(ctx, "GetRoles: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when attempting to read response: %v", err)
	}

	var result models.Roles
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal permission response to model: %v", err)
	}

	return &result, nil
}

func (c *APIClient) GetRole(ctx context.Context, id string) (*models.Roles, error) {
	uri := fmt.Sprintf(getRoleEndpoint, c.host, id)

	log.Info(ctx, "GetRole: starting permissions get policy request", log.Data{"uri": uri})

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Info(ctx, "GetRole: error building new request", log.Data{"err": err.Error()})
		return nil, err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "GetRole: error executing request", log.Data{"err": err.Error()})
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status returned from the permissions api permissions-getrole endpoint: %s", resp.Status)
	}

	log.Info(ctx, "GetRole: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when attempting to read response: %v", err)
	}

	var result models.Roles
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal permission response to model: %v", err)
	}

	return &result, nil
}

// == Policies Endpoint ==

func (c *APIClient) PostPolicy(ctx context.Context, policy models.PolicyInfo) (*models.Policy, error) {
	uri := fmt.Sprintf(addPolicyEndpoint, c.host)

	log.Info(ctx, "PostPolicy: starting permissions add policies request", log.Data{"uri": uri})

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(policy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, &buf)
	if err != nil {
		log.Info(ctx, "PostPolicy: error building new request", log.Data{"err": err.Error()})
		return nil, err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "PostPolicy: error executing request", log.Data{"err": err.Error()})
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status returned from the permissions api permissions-addpolicy endpoint: %s", resp.Status)
	}

	log.Info(ctx, "PostPolicy: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when attempting to read response: %v", err)
	}

	var result models.Policy
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal permission response to model: %v", err)
	}

	return &result, nil
}

func (c *APIClient) DeletePolicy(ctx context.Context, id string) error {
	uri := fmt.Sprintf(policyEndpoint, c.host, id)

	log.Info(ctx, "DeletePolicy: starting permissions delete policies request", log.Data{"uri": uri})

	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		log.Info(ctx, "DeletePolicy: error building new request", log.Data{"err": err.Error()})
		return err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "DeletePolicy: error executing request", log.Data{"err": err.Error()})
		return err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	log.Info(ctx, "DeletePolicy: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status returned from the permissions api permissions-deletepolicy endpoint: %s", resp.Status)
	}

	return nil
}

func (c *APIClient) GetPolicy(ctx context.Context, id string) (*models.Policy, error) {
	uri := fmt.Sprintf(policyEndpoint, c.host, id)

	log.Info(ctx, "GetPolicy: starting permissions get policy request", log.Data{"uri": uri})

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Info(ctx, "GetPolicy: error building new request", log.Data{"err": err.Error()})
		return nil, err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "GetPolicy: error executing request", log.Data{"err": err.Error()})
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status returned from the permissions api permissions-getpolicy endpoint: %s", resp.Status)
	}

	log.Info(ctx, "GetPolicy: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when attempting to read response: %v", err)
	}

	var result models.Policy
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal permission response to model: %v", err)
	}

	return &result, nil
}

func (c *APIClient) PutPolicy(ctx context.Context, id string, policy models.Policy) error {
	uri := fmt.Sprintf(policyEndpoint, c.host, id)

	log.Info(ctx, "PutPolicy: starting permissions get policy request", log.Data{"uri": uri})

	b, err := json.Marshal(policy)
	if err != nil {
		log.Info(ctx, "PutPolicy: error executing request", log.Data{"err": err.Error()})
		return err
	}

	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewReader(b))
	if err != nil {
		log.Info(ctx, "PutPolicy: error building new request", log.Data{"err": err.Error()})
		return err
	}

	resp, err := c.httpCli.Do(ctx, req)
	if err != nil {
		log.Info(ctx, "PutPolicy: error executing request", log.Data{"err": err.Error()})
		return err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status returned from the permissions api permissions-putpolicy endpoint: %s", resp.Status)
	}

	log.Info(ctx, "PutPolicy: request successfully executed", log.Data{"resp.StatusCode": resp.StatusCode})

	return nil
}

// == Permissions Endpoint ==

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
