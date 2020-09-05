package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opencars/alpr/pkg/objectstore/mockobjstore"
	"github.com/opencars/alpr/pkg/recognizer"
	"github.com/opencars/alpr/pkg/recognizer/mockalpr"
	"github.com/opencars/alpr/pkg/store"
	"github.com/opencars/alpr/pkg/store/mockstore"
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

	recognition := store.Recognition{
		ImageKey: "plates/62ae874c38d8e208f953ef90965c11f6.jpeg",
		Plate:    "AA9008MT",
	}

	mockALPR := mockalpr.NewMockRecognizer(ctrl)
	mockALPR.EXPECT().Recognize(gomock.Any()).Return(res, nil)

	mockObjStore := mockobjstore.NewMockObjectStore(ctrl)
	mockObjStore.EXPECT().Put(gomock.Any(), recognition.ImageKey, gomock.Any()).Return(nil)

	recRepo := mockstore.NewMockRecognitionRepository(ctrl)
	recRepo.EXPECT().Create(gomock.Any(), &recognition).Return(nil)

	mockStore := mockstore.NewMockStore(ctrl)
	mockStore.EXPECT().Recognition().Return(recRepo)

	srv := newServer(mockALPR, mockObjStore, mockStore)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test/example.jpeg")
	}))

	path := fmt.Sprintf("/api/v1/alpr/private/recognize?image_url=%s", s.URL)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "[{\"coordinates\":[{\"x\":431,\"y\":493},{\"x\":699,\"y\":493},{\"x\":699,\"y\":546},{\"x\":431,\"y\":546}],\"plate\":\"AA9008MT\"}]\n", rr.Body.String())
}
