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

func Test_internalServerError(t *testing.T) {
	res := internalServerError(errorMessage)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = internalServerError("")
	assert.NotEmpty(t, res.Error())
}

func Test_notFoundError(t *testing.T) {
	res := notFoundError(errorMessage)
	assert.Equal(t, http.StatusNotFound, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = notFoundError("")
	assert.NotEmpty(t, res.Error())
}

func Test_unauthorizedError(t *testing.T) {
	res := unauthorizedError(errorMessage)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = unauthorizedError("")
	assert.NotEmpty(t, res.Error())
}

func Test_forbiddenError(t *testing.T) {
	res := forbiddenError(errorMessage)
	assert.Equal(t, http.StatusForbidden, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = forbiddenError("")
	assert.NotEmpty(t, res.Error())
}

func Test_badRequestError(t *testing.T) {
	res := badRequestError(errorMessage)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Equal(t, errorMessage, res.Error())
	res = badRequestError("")
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
