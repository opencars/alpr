package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/opencars/httputil"
	"github.com/opencars/seedwork/logger"

	"github.com/opencars/alpr/pkg/domain"
	"github.com/opencars/alpr/pkg/domain/model"
	"github.com/opencars/alpr/pkg/objectstore"
	"github.com/opencars/alpr/pkg/queue"
	"github.com/opencars/alpr/pkg/store"
)

const (
	// MaxImageSize equals to 5 MB.
	MaxImageSize = 5 << 20

	// ClientTimeOut equals to 10 seconds.
	ClientTimeOut = 10 * time.Second
)

type server struct {
	router     *mux.Router
	client     *http.Client
	recognizer domain.Recognizer
	obj        objectstore.ObjectStore
	store      store.Store
	pub        queue.Publisher
}

func newServer(rec domain.Recognizer, pub queue.Publisher) *server {
	httpClient := http.Client{
		Timeout: ClientTimeOut,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ClientTimeOut,
			}).DialContext,
		},
	}

	s := server{
		router:     mux.NewRouter(),
		recognizer: rec,
		client:     &httpClient,
		pub:        pub,
	}

	s.configureRouter()

	return &s
}

func (s *server) Recognize() httputil.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		imageURL := r.URL.Query().Get("image_url")
		if imageURL == "" {
			return model.ErrRequiredImageURL
		}

		_, err := url.ParseRequestURI(imageURL)
		if err != nil {
			return model.ErrInvalidImageURL
		}

		resp, err := s.client.Get(imageURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyWithLimit := io.LimitReader(resp.Body, MaxImageSize+1)

		var buff bytes.Buffer
		if _, err = io.CopyN(&buff, bodyWithLimit, bytes.MinRead); err != nil {
			return err
		}

		typ := http.DetectContentType(buff.Bytes())
		if typ != "image/jpeg" {
			return model.ErrUnknownContentType
		}

		_, err = buff.ReadFrom(bodyWithLimit)
		if err != nil {
			return err
		}

		if buff.Len() > MaxImageSize {
			return model.ErrImageTooLarge
		}

		reader := bytes.NewReader(buff.Bytes())

		res, err := s.recognizer.Recognize(reader)
		if err != nil {
			return err
		}

		if len(res) > 0 {
			err := s.pub.Publish(&queue.Event{
				URL:    imageURL,
				Number: res[0].Plate,
			})

			if err != nil {
				logger.Errorf("publish: %v", err)
			}
		}

		return json.NewEncoder(w).Encode(res)
	}
}
