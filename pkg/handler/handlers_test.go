package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Run("StatusError", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		h := Handler(func(w http.ResponseWriter, r *http.Request) error {
			return NewError(http.StatusBadRequest, "blah-blah")
		})

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"code\":400,\"message\":\"blah-blah\"}\n", rr.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		h := Handler(func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("blah-blah")
		})

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "{\"code\":500,\"message\":\"Something went wrong\"}\n", rr.Body.String())
	})
}
