package apierrors

import (
	"errors"
)

//A lit of error messages for Permissions API
var (
	//ErrRoleNotFound is an error when the role can not be found in mongoDB
	ErrRoleNotFound = errors.New("role not found")
)
