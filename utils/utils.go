package utils

import (
	"strconv"

	"github.com/ONSdigital/dp-permissions-api/apierrors"
)

//ValidatePositiveInteger checks if a value is positive
func ValidatePositiveInteger(value string) (int, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return -1, apierrors.ErrInvalidPositiveInteger
	}
	if val < 0 {
		return -1, apierrors.ErrInvalidPositiveInteger
	}
	return val, nil
}
