package sdk

import "errors"

var (

	// ErrGetPermissionsResponseBodyNil error used when a nil response is returned from the permissions API.
	ErrGetPermissionsResponseBodyNil = errors.New("error creating get permissions request http.Request required but was nil")

	// ErrFailedToParsePermissionsResponse error used when an unexpected response body is returned from the permissions API and it fails to parse.
	ErrFailedToParsePermissionsResponse = errors.New("error parsing permissions bundle response body")

	// ErrNotCached error error used when permissions are not found in the cache.
	ErrNotCached = errors.New("permissions bundle not found in the cache")
)
