package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPError(t *testing.T) {
	t.Run("Test HTTP Error",
		func(t *testing.T) {
			err := HTTPError(http.StatusBadRequest, "Some Error")
			assert.Equal(
				t,
				"code=400, message={400 Some Error}",
				err.Error(),
			)
		})
}
