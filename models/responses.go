package models

// ErrorResponse represents a slice of errors in a JSON response body.
type ErrorResponse struct {
	Errors  []error           `json:"errors"`
	Status  int               `json:"-"`
	Headers map[string]string `json:"-"`
}

func NewErrorResponse(statusCode int, headers map[string]string, errors ...error) *ErrorResponse {
	return &ErrorResponse{
		Errors:  errors,
		Status:  statusCode,
		Headers: headers,
	}
}

// SuccessResponse represents a success JSON response body.
type SuccessResponse struct {
	Body    []byte            `json:"-"`
	Status  int               `json:"-"`
	Headers map[string]string `json:"-"`
}

// NewSuccessResponse creates a new SuccessResponse.
func NewSuccessResponse(jsonBody []byte, statusCode int, headers map[string]string) *SuccessResponse {
	return &SuccessResponse{
		Body:    jsonBody,
		Status:  statusCode,
		Headers: headers,
	}
}
