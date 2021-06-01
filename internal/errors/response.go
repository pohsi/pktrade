package errors

import (
	"net/http"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
)

type errorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e errorResponse) Error() string {
	return e.Message
}

func (e errorResponse) StatusCode() int {
	return e.Status
}

func InternalServerError(msg string) errorResponse {

	if msg == "" {
		msg = "Encountered an unknow error"
	}

	return errorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

func NotFoundError(msg string) errorResponse {

	if msg == "" {
		msg = "The requested resource was not found"
	}

	return errorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

func UnauthorizedError(msg string) errorResponse {

	if msg == "" {
		msg = "Not authennicated to perfoem request action"
	}
	return errorResponse{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}

func ForbiddenError(msg string) errorResponse {

	if msg == "" {
		msg = "Not authorized to perfoem request action"
	}

	return errorResponse{
		Status:  http.StatusForbidden,
		Message: msg,
	}
}

func BadRequestError(msg string) errorResponse {

	if msg == "" {
		msg = "Bad request"
	}

	return errorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

type invalidField struct {
	Field string `json:"field"`
	Error string `jsom:"error"`
}

func invalidInput(errs validation.Errors) errorResponse {
	var details []invalidField
	var fields []string
	for field := range errs {
		fields = append(fields, field)
	}

	sort.Strings(fields)
	for _, field := range fields {
		details = append(details, invalidField{
			Field: field,
			Error: errs[field].Error(),
		})
	}

	return errorResponse{
		Status:  http.StatusBadRequest,
		Message: "Invalid input",
		Details: details,
	}
}
