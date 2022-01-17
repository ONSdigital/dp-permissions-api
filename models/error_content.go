package models

// API error codes
const (
	JSONMarshalError                           = "JSONMarshalError"
	JSONUnmarshalError                         = "JSONUnmarshalError"
	WriteResponseError                         = "WriteResponseError"
	InvalidQueryParameterError                 = "InvalidQueryParameter" // TODO: create query parameter specific errors?
	InvalidLimitQueryParameterMaxExceededError = "InvalidLimitQueryParameterMaxExceeded"
	RoleNotFoundError                          = "RoleNotFoundError"
	GetRoleError                               = "GetRoleError"
	GetRolesError                              = "GetRolesError"
	GetPermissionBundleError                   = "GetPermissionBundleError"
	PolicyNotFoundError                        = "PolicyNotFoundError"
	GetPolicyError                             = "GetPolicyError"
	DeletePolicyError                          = "DeletePolicyError"
	InvalidPolicyError                         = "InvalidPolicyError"
	CreateNewPolicyError                       = "CreateNewPolicyError"
	UpdatePolicyError                          = "UpdatePolicyError"
)

// API error descriptions
const (
	MarshalFailedDescription                         = "failed to marshal the request body"
	UnmarshalFailedDescription                       = "unable to unmarshal request body"
	ErrorMarshalFailedDescription                    = "failed to marshal the error"
	WriteResponseFailedDescription                   = "failed to write http response"
	InvalidQueryParameterDescription                 = "invalid query parameter: "
	InvalidLimitQueryParameterMaxExceededDescription = "invalid query parameter: limit, maximum limit exceeded"
	RoleNotFoundDescription                          = "role not found"
	GetRoleErrorDescription                          = "retrieving role from DB returned an error"
	GetRolesErrorDescription                         = "retrieving roles from DB returned an error"
	GetPermissionBundleErrorDescription              = "failed to get permissions bundle"
	PolicyNotFoundDescription                        = "policy not found"
	GetPolicyErrorDescription                        = "retrieving policy from DB returned an error"
	DeletePolicyErrorDescription                     = "deleting policy from DB returned an error"
	InvalidPolicyDescription                         = "policy parameters failed validation: "
	CreateNewPolicyErrorDescription                  = "failed to create new policy"
	UpdatePolicyErrorDescription                     = "failed to update policy"
)
