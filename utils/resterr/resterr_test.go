package resterr_test

import (
	"net/http"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
	"github.com/stretchr/testify/require"
)

func TestRestErr_Error(t *testing.T) {
	testCases := []struct {
		name          string
		function      func(string) *resterr.RestErr
		inputMessage  string
		expectedError *resterr.RestErr
	}{
		{
			name:          "Test_NewBadRequestError",
			function:      resterr.NewBadRequestError,
			inputMessage:  "bad request",
			expectedError: &resterr.RestErr{Message: "bad request", Err: "bad_request", Code: http.StatusBadRequest},
		},
		{
			name:          "Test_NewUnauthorizedRequestError",
			function:      resterr.NewUnauthorizedRequestError,
			inputMessage:  "Unauthorized",
			expectedError: &resterr.RestErr{Message: "Unauthorized", Err: "unauthorized", Code: http.StatusUnauthorized},
		},
		{
			name:          "Test_NewNotFoundError",
			function:      resterr.NewNotFoundError,
			inputMessage:  "not found",
			expectedError: &resterr.RestErr{Message: "not found", Err: "not_found", Code: http.StatusNotFound},
		},
		{
			name:          "Test_NewInternalServerError",
			function:      resterr.NewInternalServerError,
			inputMessage:  "internal server error",
			expectedError: &resterr.RestErr{Message: "internal server error", Err: "internal_server_error", Code: http.StatusInternalServerError},
		},
		{
			name:          "Test_NewForbiddenError",
			function:      resterr.NewForbiddenError,
			inputMessage:  "forbidden",
			expectedError: &resterr.RestErr{Message: "forbidden", Err: "forbidden", Code: http.StatusForbidden},
		},
		{
			name:          "Test_NewConflictError",
			function:      resterr.NewConflictError,
			inputMessage:  "conflict",
			expectedError: &resterr.RestErr{Message: "conflict", Err: "conflict", Code: http.StatusConflict},
		},
		{
			name:          "Test_NewServiceUnavailableError",
			function:      resterr.NewServiceUnavailableError,
			inputMessage:  "service unavailable",
			expectedError: &resterr.RestErr{Message: "service unavailable", Err: "service_unavailable", Code: http.StatusServiceUnavailable},
		},
		{
			name:          "Test_NewTooManyRequests",
			function:      resterr.NewTooManyRequestsError,
			inputMessage:  "too many requests",
			expectedError: &resterr.RestErr{Message: "too many requests", Err: "too_many_requests", Code: http.StatusTooManyRequests},
		},
		{
			name:          "Test_NewUnprocessableEntityError",
			function:      resterr.NewUnprocessableEntityError,
			inputMessage:  "unprocessable entity",
			expectedError: &resterr.RestErr{Message: "unprocessable entity", Err: "unprocessable_entity", Code: http.StatusUnprocessableEntity},
		},
		{
			name:          "Test_NewBadGatewayError",
			function:      resterr.NewBadGatewayError,
			inputMessage:  "bad gateway",
			expectedError: &resterr.RestErr{Message: "bad gateway", Err: "bad_gateway", Code: http.StatusBadGateway},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.function(tc.inputMessage)
			require.Equal(t, tc.expectedError, result)
		})
	}
}

func TestRestErrWithCauses_Error(t *testing.T) {
	testCases := []struct {
		name          string
		function      func(string, []resterr.Causes) *resterr.RestErr
		inputMessage  string
		inputCauses   []resterr.Causes
		expectedError *resterr.RestErr
	}{
		{
			name:         "Test_NewBadRequestValidationError",
			function:     resterr.NewBadRequestValidationError,
			inputMessage: "bad request",
			inputCauses: []resterr.Causes{
				{
					Field:   "error",
					Message: "error",
				},
			},
			expectedError: &resterr.RestErr{
				Message: "bad request",
				Err:     "bad_request",
				Code:    http.StatusBadRequest,
				Causes: []resterr.Causes{
					{
						Field:   "error",
						Message: "error",
					},
				},
			},
		},
		{
			name:         "Test_NewUnprocessableEntityWithCausesError",
			function:     resterr.NewUnprocessableEntityWithCausesError,
			inputMessage: "unprocessable entity",
			inputCauses: []resterr.Causes{
				{
					Field:   "error",
					Message: "error",
				},
			},
			expectedError: &resterr.RestErr{
				Message: "unprocessable entity",
				Err:     "unprocessable_entity",
				Code:    http.StatusUnprocessableEntity,
				Causes: []resterr.Causes{
					{
						Field:   "error",
						Message: "error",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.function(tc.inputMessage, tc.inputCauses)
			require.Equal(t, tc.expectedError, result)
		})
	}
}

func TestRestErrError_Error(t *testing.T) {
	testCases := []struct {
		expectedMessage string
		err             *resterr.RestErr
	}{
		{
			expectedMessage: "error",
			err:             &resterr.RestErr{Message: "error", Err: "error", Code: http.StatusInternalServerError},
		},
		{
			expectedMessage: "test",
			err:             &resterr.RestErr{Message: "test", Err: "test", Code: http.StatusInternalServerError},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedMessage, func(t *testing.T) {
			require.Equal(t, tc.err.Error(), tc.err.Message)
		})
	}
}
