package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/opencars/alpr/pkg/handler"
	"github.com/opencars/alpr/pkg/logger"
	"github.com/opencars/alpr/pkg/objectstore"
	"github.com/opencars/alpr/pkg/recognizer"
	"github.com/opencars/alpr/pkg/version"
)

const (
	// MaxImageSize equals to 5 MB.
	MaxImageSize = 5 << 20

	// ClientTimeOut equals to 10 seconds.
	ClientTimeOut = 10 * time.Second
)

type server struct {
	router      *mux.Router
	client      *http.Client
	recognizer  recognizer.Recognizer
	objectStore objectstore.ObjectStore
}

func newServer(recognizer recognizer.Recognizer, objectStore objectstore.ObjectStore) *server {
	httpClient := http.Client{
		Timeout: ClientTimeOut,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ClientTimeOut,
			}).DialContext,
		},
	}

	s := server{
		router:      mux.NewRouter(),
		recognizer:  recognizer,
		objectStore: objectStore,
		client:      &httpClient,
	}

	s.configureRouter()

	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Api-Key", "X-Api-Key"})

	cors := handlers.CORS(origins, methods, headers)(s.router)
	compress := handlers.CompressHandler(cors)
	compress.ServeHTTP(w, r)
}

func (*server) Swagger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yml")
	}
}

func (*server) Version() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		v := struct {
			Version string `json:"version"`
			Go      string `json:"go"`
		}{
			Version: version.Version,
			Go:      runtime.Version(),
		}

		return json.NewEncoder(w).Encode(v)
	}
}

// Note: Later we could publish event of recognition into the NATS queue and prepare.
func (s *server) Recognize() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		imageURL := r.URL.Query().Get("image_url")
		if imageURL == "" {
			return handler.ErrRequiredImageURL
		}

		_, err := url.ParseRequestURI(imageURL)
		if err != nil {
			return handler.ErrInvalidImageURL
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
			return handler.ErrUnknownContentType
		}

		_, err = buff.ReadFrom(bodyWithLimit)
		if err != nil {
			return err
		}

		if buff.Len() > MaxImageSize {
			return handler.ErrImageTooLarge
		}

		reader := bytes.NewReader(buff.Bytes())

		res, err := s.recognizer.Recognize(reader)
		if err != nil {
			return err
		}

		if len(res) > 0 && s.objectStore != nil {
			// Reset the read pointer.
			_, err = reader.Seek(0, 0)
			if err != nil {
				return err
			}

			err := s.objectStore.Put(r.Context(), reader)
			if err != nil {
				logger.Errorf("failed to put: %v", err)
			}

			// TODO: Save number and URL to store!
		}

		return json.NewEncoder(w).Encode(res)
	}
}
