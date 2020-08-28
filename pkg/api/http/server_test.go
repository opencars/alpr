package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opencars/alpr/pkg/recognizer/mockalpr"
)

func TestServer_Recognize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// TODO: Add test result to model.

	recognizer := mockalpr.NewMockRecognizer(ctrl)
	recognizer.EXPECT().Recognize(gomock.Any()).Return(nil, nil)

	srv := newServer(recognizer, nil)

	imageURL := "https://example.com/vehicle.jpg"
	path := fmt.Sprintf("/api/v1/alpr/private/recognize?image_url=%s", imageURL)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
