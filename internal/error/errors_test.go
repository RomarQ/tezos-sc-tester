package error

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPError(t *testing.T) {
	t.Run("Test Basic HTTP Error",
		func(t *testing.T) {
			err := HttpError(http.StatusBadRequest, "Some Error")
			assert.Equal(
				t,
				"code=400, message={400 Some Error <nil>}",
				err.Error(),
			)
		})
	t.Run("Test Detailed HTTP Error",
		func(t *testing.T) {
			err := DetailedHttpError(http.StatusBadRequest, "Some Error", [...]interface{}{})
			assert.Equal(
				t,
				"code=400, message={400 Some Error []}",
				err.Error(),
			)
		})
}
