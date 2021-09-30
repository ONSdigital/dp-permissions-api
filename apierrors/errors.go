package apierrors

import (
	"errors"
	"strconv"
)

//A lit of error messages for Permissions API
var (
	//ErrRoleNotFound is an error when the role can not be found in mongoDB
	ErrRoleNotFound           = errors.New("role not found")
	ErrInvalidPositiveInteger = errors.New("value is not a positive integer")
	ErrLimitAndOffset         = errors.New("offset and limit must be positive or zero")
	ErrPolicyNotFound           = errors.New("policy not found")
)

// ErrorMaximumLimitReached creates a unique error
func ErrorMaximumLimitReached(m int) error {
	err := errors.New("the maximum limit has been reached, the limit cannot be more than " + strconv.Itoa(m))
	return err
}
