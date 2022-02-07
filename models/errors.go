package models

import (
	"context"

	"github.com/ONSdigital/log.go/v2/log"
)

// Error represents an error.
type Error struct {
	Cause       error  `json:"-"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Error returns a string representation of the error. Implements error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Code + ": " + e.Description
}

// NewError creates a new Error. Once created. the error is logged along with logData.
func NewError(ctx context.Context, cause error, code string, description string, logData log.Data) *Error {
	err := &Error{
		Cause:       cause,
		Code:        code,
		Description: description,
	}
	log.Error(ctx, description, err, logData)
	return err
}
