package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opencars/alpr/pkg/queue"
	"github.com/opencars/alpr/pkg/queue/mockqueue"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opencars/alpr/pkg/recognizer"
	"github.com/opencars/alpr/pkg/recognizer/mockalpr"
)

func TestServer_Recognize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	res := []recognizer.Result{
		{
			Coordinates: []recognizer.Coordinate{
				{
					X: 431,
					Y: 493,
				},
				{
					X: 699,
					Y: 493,
				},
				{
					X: 699,
					Y: 546,
				},
				{
					X: 431,
					Y: 546,
				},
			},
			Plate: "AA9008MT",
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test/example.jpeg")
	}))

	event := queue.Event{
		URL:    s.URL,
		Number: res[0].Plate,
	}

	mockALPR := mockalpr.NewMockRecognizer(ctrl)
	mockALPR.EXPECT().Recognize(gomock.Any()).Return(res, nil)

	mockPublisher := mockqueue.NewMockPublisher(ctrl)
	mockPublisher.EXPECT().Publish(&event)

	srv := newServer(mockALPR, mockPublisher)

	path := fmt.Sprintf("/api/v1/alpr/private/recognize?image_url=%s", s.URL)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "[{\"coordinates\":[{\"x\":431,\"y\":493},{\"x\":699,\"y\":493},{\"x\":699,\"y\":546},{\"x\":431,\"y\":546}],\"plate\":\"AA9008MT\"}]\n", rr.Body.String())
}
