package models

import (
	"context"

	"github.com/ONSdigital/log.go/v2/log"
)

type Error struct {
	Cause       error  `json:"-"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Code + ": " + e.Description
}

func NewError(ctx context.Context, cause error, code string, description string, logData log.Data) *Error {
	err := &Error{
		Cause:       cause,
		Code:        code,
		Description: description,
	}
	log.Error(ctx, description, err, logData)
	return err
}
