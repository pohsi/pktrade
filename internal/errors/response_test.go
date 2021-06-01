package errors

import (
	"fmt"
	"net/http"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

const errorMessage = "test"

func Test_errorResponce(t *testing.T) {
	const statusCode = 400

	e := errorResponse{
		Message: errorMessage,
		Status:  statusCode,
	}
	assert.Equal(t, errorMessage, e.Error())
	assert.Equal(t, statusCode, e.StatusCode())
}

func TestInternalServerError(t *testing.T) {
	res := InternalServerError(errorMessage)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = InternalServerError("")
	assert.NotEmpty(t, res.Error())
}

func TestNotFoundError(t *testing.T) {
	res := NotFoundError(errorMessage)
	assert.Equal(t, http.StatusNotFound, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = NotFoundError("")
	assert.NotEmpty(t, res.Error())
}

func TestUnauthorizedError(t *testing.T) {
	res := UnauthorizedError(errorMessage)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = UnauthorizedError("")
	assert.NotEmpty(t, res.Error())
}

func Test_forbiddenError(t *testing.T) {
	res := ForbiddenError(errorMessage)
	assert.Equal(t, http.StatusForbidden, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = ForbiddenError("")
	assert.NotEmpty(t, res.Error())
}

func TestBadRequestError(t *testing.T) {
	res := BadRequestError(errorMessage)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = BadRequestError("")
	assert.NotEmpty(t, res.Error())
}
func Test_invalidInput(t *testing.T) {
	err := invalidInput(validation.Errors{
		"xyz": fmt.Errorf("2"),
		"abc": fmt.Errorf("1"),
	})

	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, []invalidField{{"abc", "1"}, {"xyz", "2"}}, err.Details)
}
