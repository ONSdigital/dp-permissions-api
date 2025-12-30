package models

// API error codes
const (
	JSONMarshalError                           = "JSONMarshalError"
	JSONUnmarshalError                         = "JSONUnmarshalError"
	WriteResponseError                         = "WriteResponseError"
	InvalidQueryParameterError                 = "InvalidQueryParameter"
	InvalidLimitQueryParameterMaxExceededError = "InvalidLimitQueryParameterMaxExceeded"
	RoleNotFoundError                          = "RoleNotFoundError"
	GetRoleError                               = "GetRoleError"
	GetRolesError                              = "GetRolesError"
	GetPermissionBundleError                   = "GetPermissionBundleError"
	PolicyNotFoundError                        = "PolicyNotFoundError"
	PolicyAlreadyExistsError                   = "PolicyAlreadyExistsError"
	GetPolicyError                             = "GetPolicyError"
	DeletePolicyError                          = "DeletePolicyError"
	InvalidPolicyError                         = "InvalidPolicyError"
	CreateNewPolicyError                       = "CreateNewPolicyError"
	CreatePolicyWithIDError                    = "CreatePolicyWithIDError"
	UpdatePolicyError                          = "UpdatePolicyError"
)

// API error descriptions
const (
	InternalServerErrorDescription                   = "internal server error"
	MarshalFailedDescription                         = "failed to marshal the request body"
	UnmarshalFailedDescription                       = "unable to unmarshal request body"
	ErrorMarshalFailedDescription                    = "failed to marshal the error"
	WriteResponseFailedDescription                   = "failed to write http response"
	InvalidQueryParameterDescription                 = "invalid query parameter"
	InvalidLimitQueryParameterMaxExceededDescription = "invalid query parameter: maximum exceeded"
	RoleNotFoundDescription                          = "role not found"
	GetRoleErrorDescription                          = "retrieving role from DB returned an error"
	GetRolesErrorDescription                         = "retrieving roles from DB returned an error"
	GetPermissionBundleErrorDescription              = "failed to get permissions bundle"
	PolicyAlreadyExistsDescription                   = "policy already exists with given ID"
	PolicyNotFoundDescription                        = "policy not found"
	GetPolicyErrorDescription                        = "retrieving policy from DB returned an error"
	DeletePolicyErrorDescription                     = "deleting policy from DB returned an error"
	CreateNewPolicyErrorDescription                  = "failed to create new policy"
	CreatePolicyWithIDErrorDescription               = "failed to create policy with given ID"
	UpdatePolicyErrorDescription                     = "failed to update policy"
)
