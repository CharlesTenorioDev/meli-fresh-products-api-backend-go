package resterr

import "net/http"

// RestErr represents the error object.
// @Summary Error information
// @Description Structure for describing why the error occurred
type RestErr struct {
	// Error message.
	Message string `json:"message" example:"error trying to process request"`

	// Error description.
	Err string `json:"error" example:"internal_server_error"`

	// Error code.
	Code int `json:"code" example:"500"`

	// Error causes.
	Causes []Causes `json:"causes"`
}

// Causes represents the error causes.
// @Summary Error Causes
// @Description Structure representing the causes of an error.
type Causes struct {
	// Field associated with the error cause.
	// @json
	// @jsonTag field
	Field string `json:"field" example:"name"`

	// Error message describing the cause.
	// @json
	// @jsonTag message
	Message string `json:"message" example:"name is required"`
}

// Error implements the error interface for RestErr.
// It returns the error message contained within the RestErr object.
func (r *RestErr) Error() string {
	return r.Message
}

// NewBadRequestError returns a RestErr with http.StatusBadRequest code.
//
// It should be used when the client sends a request that the server cannot or will not process
// due to a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing).
func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

// NewUnauthorizedRequestError returns a RestErr with http.StatusUnauthorized code.
//
// It should be used when the request requires user authentication, and the client has not
// provided valid authentication credentials.
func NewUnauthorizedRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}

// NewBadRequestValidationError returns a RestErr with http.StatusBadRequest code and causes.
//
// It should be used when the request is well-formed, but the server is unable to process the contained instructions.
// For example, this can be used when the request body contains invalid data.
func NewBadRequestValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

// NewInternalServerError returns a RestErr with http.StatusInternalServerError code.
//
// It should be used when the server encountered an unexpected condition that prevented it from fulfilling the request.
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

// NewNotFoundError returns a RestErr with http.StatusNotFound code.
//
// It should be used when the server could not find the requested resource.
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	}
}

// NewForbiddenError returns a RestErr with http.StatusForbidden code.
//
// It should be used when the server understood the request, but is
// refusing to authorize it.
func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}

// NewConflictError returns a RestErr with http.StatusConflict code.
//
// It should be used when the request conflicts with an existing resource.
func NewConflictError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "conflict",
		Code:    http.StatusConflict,
	}
}

// NewUnprocessableEntityError returns a RestErr with http.StatusUnprocessableEntity code.
//
// It should be used when the request is well-formed, but the server is unable to process the contained instructions.
func NewUnprocessableEntityError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unprocessable_entity",
		Code:    http.StatusUnprocessableEntity,
	}
}

// NewBadGatewayError returns a RestErr with http.StatusBadGateway code.
//
// It should be used when the server, while acting as a gateway or proxy, received an invalid response from the upstream server.
func NewBadGatewayError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_gateway",
		Code:    http.StatusBadGateway,
	}
}

// NewServiceUnavailableError returns a RestErr with http.StatusServiceUnavailable code.
//
// It should be used when the service is currently unable to handle the request due to a temporary overloading or maintenance of the server.
func NewServiceUnavailableError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "service_unavailable",
		Code:    http.StatusServiceUnavailable,
	}
}

// NewTooManyRequestsError returns a RestErr with http.StatusTooManyRequests code.
//
// It should be used when the client has sent too many requests in a given amount of time.
func NewTooManyRequestsError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "too_many_requests",
		Code:    http.StatusTooManyRequests,
	}
}

// NewUnprocessableEntityWithCausesError returns a RestErr with http.StatusUnprocessableEntity code and causes.
//
// It should be used when the request is well-formed, but the server is unable to process the contained instructions.
func NewUnprocessableEntityWithCausesError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unprocessable_entity",
		Code:    http.StatusUnprocessableEntity,
		Causes:  causes,
	}
}
